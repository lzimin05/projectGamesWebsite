package api

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-playground/validator"
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

func (srv *Server) Login(e echo.Context) error {
	input := struct {
		Email    string `json:"email" validate:"required,email"`
		Password string `json:"password" validate:"required"`
	}{}
	err := e.Bind(&input)
	if err != nil {
		return e.String(http.StatusInternalServerError, err.Error())
	}
	password, err := srv.uc.GetPasswordByEmail(input.Email)
	if err != nil {
		if errors.Is(err, vars.ErrAlreadyExist) {
			return e.String(http.StatusConflict, err.Error())
		}
		return e.String(http.StatusInternalServerError, err.Error())
	}
	validate := validator.New()
	err = validate.Struct(input)
	if err != nil {
		return e.String(http.StatusBadRequest, "почта указана неверно")
	}
	if password == "" {
		fmt.Println("Пользователь еще не зарегестрирован!")
		return e.String(http.StatusBadRequest, "Пользователь еще не зарегестрирован!")
	}
	flag := newuser.CheckPasswordHash(input.Password, password)
	if flag {
		fmt.Println("совпал")
	} else {
		fmt.Println("не совпал пароль")
	}
	return e.JSON(http.StatusOK, flag)
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
