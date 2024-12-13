package api

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	maxSize int

	server  *echo.Echo
	address string

	uc Usecase
}

func NewServer(ip string, port int, maxSize int, uc Usecase) *Server {
	api := Server{
		maxSize: maxSize,
		uc:      uc,
	}

	api.server = echo.New()

	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                                          // Разрешаем запросы с вашего frontend адреса
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodOptions},          // Разрешенные методы
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept}, // Разрешенные заголовки
	}))

	api.server.GET("/user/email", api.GetUserByEmail)
	api.server.GET("/user/:id", api.GetUserById)
	api.server.POST("/newuser", api.PostNewUser)

	api.server.POST("/login", api.Login)

	api.address = fmt.Sprintf("%s:%d", ip, port)

	return &api
}

func (api *Server) Run() {
	api.server.Logger.Fatal(api.server.Start(api.address))
}
