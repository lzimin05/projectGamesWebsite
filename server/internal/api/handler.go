package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	newuser "github.com/lzimin05/IDZ/internal/user"
	"github.com/lzimin05/IDZ/pkg/vars"
)

func (srv *Server) GetUserById(e echo.Context) error {
	idparam, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid id")
	}
	msg, err := srv.uc.PrintUserById(idparam)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, msg)
}

func (srv *Server) GetEmailById(e echo.Context) error {
	idparam, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusBadRequest, "invalid id")
	}
	msg, err := srv.uc.PrintEmailById(idparam)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, msg)
}

func (srv *Server) GetUserByEmail(e echo.Context) error {
	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	email := claims.Email
	if email == "" {
		return e.String(http.StatusBadRequest, "email is empty")
	}
	msg, err := srv.uc.PrintUserByEmail(email)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, msg)
}

func (srv *Server) Access(e echo.Context) error {
	user := e.Get("user").(*jwt.Token)
	claims := user.Claims.(*jwtCustomClaims)
	email := claims.Email
	return e.String(http.StatusOK, fmt.Sprintf("Добро пожаловать, %s!", email))
}

func (srv *Server) Login(e echo.Context) error {
	input := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	err := e.Bind(&input)
	if err != nil {
		return e.String(http.StatusBadRequest, "Неверные данные")
	}

	// Проверка пароля
	password, err := srv.uc.GetPasswordByEmail(input.Email)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Ошибка при проверке пароля")
	}

	if !newuser.CheckPasswordHash(input.Password, password) {
		return e.String(http.StatusUnauthorized, "Неверный email или пароль")
	}

	// Создаем JWT-токен
	claims := &jwtCustomClaims{
		Email: input.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)), //время жизни токена
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Подписываем токен секретным ключом
	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Не удалось создать токен")
	}

	// Добавляем сессию
	id, err := srv.uc.GetIdByemail(input.Email)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	err = srv.uc.UpdateSesionNow(id)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	// Возвращаем токен в ответе
	return e.JSON(http.StatusOK, map[string]string{
		"token": tokenString,
	})
}

func (srv *Server) PostNewUser(e echo.Context) error {
	var NewUser newuser.User
	err := e.Bind(&NewUser)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	if len([]rune(NewUser.Name)) > 30 || len([]rune(NewUser.Name)) < 5 {
		return e.String(http.StatusBadRequest, "Длина от 5 до 30")
	}

	if len([]rune(NewUser.Password)) <= 5 {
		return e.String(http.StatusBadRequest, "Длина пароля должна быть больше 5 символов")
	}
	validate := validator.New()
	err = validate.Struct(NewUser)
	if err != nil {
		return e.String(http.StatusBadRequest, "почта указана неверно")
	}
	id, err := srv.uc.GetIdByemail(NewUser.Email)
	if err != nil {
		return e.String(http.StatusBadRequest, "почта указана неверно")
	}
	if id != 0 {
		return e.String(http.StatusNoContent, "уже зарегистрирован!")
	}
	hashedpassword, err := newuser.HashPassword(NewUser.Password)
	if err != nil {
		return e.String(http.StatusInternalServerError, "Ошибка хеширования пароля")
	}
	NewUser.Password = hashedpassword
	err = srv.uc.InsertNewUser(NewUser)
	if err != nil {
		if errors.Is(err, vars.ErrAlreadyExist) {
			return e.String(http.StatusConflict, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.String(http.StatusCreated, "OK")
}

func (srv *Server) Logout(c echo.Context) error {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return c.String(http.StatusBadRequest, "Missing token")
	}

	tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

	// Добавляем токен в черный список
	revokedTokens.Lock()
	revokedTokens.tokens[tokenString] = true
	revokedTokens.Unlock()

	// Обновляем сессию на id = 0
	err := srv.uc.UpdateSesionNow(0)
	if err != nil {
		return c.String(http.StatusInternalServerError, err.Error())
	}

	return c.String(http.StatusOK, "Logged out successfully")
}

var revokedTokens = struct {
	sync.RWMutex
	tokens map[string]bool
}{tokens: make(map[string]bool)}

func (srv *Server) GetSesion(e echo.Context) error {
	id, err := srv.uc.GetSesionNow()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.String(http.StatusOK, strconv.Itoa(id))
}

func (srv *Server) UpdateUserById(e echo.Context) error {
	idparam, err := strconv.Atoi(e.Param("id"))
	if err != nil {
		return e.String(http.StatusInternalServerError, "invalid id")
	}
	idsesion, err := srv.uc.GetSesionNow()
	if err != nil {
		return e.String(http.StatusInternalServerError, "invalid id")
	}
	if idparam != idsesion {
		return e.String(http.StatusForbidden, "error access denied")
	}
	newname := struct {
		Newname string `json:"newname" validate:"required"`
	}{}
	err = e.Bind(&newname)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	if len([]rune(newname.Newname)) > 30 || len([]rune(newname.Newname)) < 5 {
		return e.String(http.StatusBadRequest, "Длина от 5 до 30")
	}
	validate := validator.New()
	err = validate.Struct(newname)
	if err != nil {
		return e.String(http.StatusBadRequest, "обновленное имя пользователя указано неверно")
	}
	err = srv.uc.UpdateUserById(newname.Newname, idparam)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.String(http.StatusOK, "OK")
}

func (srv *Server) AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Проверяем, находится ли токен в черном списке
		revokedTokens.RLock()
		if revokedTokens.tokens[tokenString] {
			_ = srv.uc.UpdateSesionNow(0)
			revokedTokens.RUnlock()
			return c.String(http.StatusUnauthorized, "Unauthorized: Token revoked")
		}
		revokedTokens.RUnlock()

		token, err := jwt.ParseWithClaims(tokenString, &jwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		// Если токен истек или недействителен, обновляем сессию на id = 0
		if err != nil {
			// Проверяем, является ли ошибка истечением токена
			if errors.Is(err, jwt.ErrTokenExpired) {
				fmt.Println("Token expired") // Логирование
				err := srv.uc.UpdateSesionNow(0)
				if err != nil {
					return c.String(http.StatusInternalServerError, err.Error())
				}
				return c.Redirect(http.StatusTemporaryRedirect, "/login")
			}
			// Если ошибка другая, просто перенаправляем на страницу входа
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
			c.Set("email", claims.Email)
			return next(c)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}
