package api

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
	newuser "github.com/lzimin05/IDZ/internal/user"
	"github.com/lzimin05/IDZ/pkg/vars"
	"golang.org/x/crypto/bcrypt"
)

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

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

func (srv *Server) GetUserByEmail(e echo.Context) error {
	input := e.FormValue("email")
	if input == "" {
		return e.String(http.StatusBadRequest, "email is empty")
	}
	msg, err := srv.uc.PrintUserByEmail(input)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, msg)
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

	if len([]rune(NewUser.Name)) <= 5 {
		return e.String(http.StatusBadRequest, "Длина пароля должна быть больше 5 символов")
	}
	validate := validator.New()
	err = validate.Struct(NewUser)
	if err != nil {
		return e.String(http.StatusBadRequest, "почта указана неверно")
	}
	hashedpassword, err := hashPassword(NewUser.Password)
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

/*
// GetHello возвращает случайное приветствие пользователю
func (srv *Server) GetHello(e echo.Context) error {
	msg, err := srv.uc.FetchHelloMessage()
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	return e.JSON(http.StatusOK, msg)
}

// PostHello Помещает новый вариант приветствия в БД
func (srv *Server) PostHello(e echo.Context) error {
	input := struct {
		Msg *string `json:"msg"`
	}{}

	err := e.Bind(&input)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}

	if input.Msg == nil {
		return e.String(http.StatusBadRequest, "msg is empty")
	}

	if len([]rune(*input.Msg)) > srv.maxSize {
		return e.String(http.StatusBadRequest, "hello message too large")
	}

	err = srv.uc.SetHelloMessage(*input.Msg)
	if err != nil {
		if errors.Is(err, vars.ErrAlreadyExist) {
			return e.String(http.StatusConflict, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}

	return e.String(http.StatusCreated, "OK")
}
*/
