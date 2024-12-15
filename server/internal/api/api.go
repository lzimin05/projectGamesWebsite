package api

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jwtSecret = []byte("your-secret-key")

type Server struct {
	maxSize int

	server  *echo.Echo
	address string

	uc Usecase
}

type jwtCustomClaims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func NewServer(ip string, port int, maxSize int, uc Usecase) *Server {
	api := Server{
		maxSize: maxSize,
		uc:      uc,
	}

	api.server = echo.New()

	// Логирование и восстановление после паники
	api.server.Use(middleware.Logger())
	api.server.Use(middleware.Recover())

	// Настройка CORS
	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://127.0.0.1:5500"},
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
	}))

	// Маршрут для входа
	api.server.POST("/newuser", api.PostNewUser)
	api.server.POST("/login", api.Login)

	// Создаем группу маршрутов /restricted
	r := api.server.Group("/restricted")

	// Middleware для проверки JWT-токена
	config := echojwt.Config{
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(jwtCustomClaims)
		},
		SigningKey: jwtSecret,
	}
	r.Use(echojwt.WithConfig(config))

	// Защищенные маршруты
	r.GET("/user/:email", api.GetUserByEmail)
	r.GET("/email/:id", api.GetEmailById)
	r.GET("/main", api.Access)
	r.GET("/:id", api.GetUserById)
	r.GET("/sesion", api.GetSesion)
	r.PUT("/user/:id", api.UpdateUserById)

	// Маршрут для страницы входа
	api.server.GET("/login", func(c echo.Context) error {
		return c.File("../../../registation/auth.html") // Возвращаем страницу входа
	})

	api.server.POST("/logout", api.Logout)

	// Маршруты для статических файлов с middleware для проверки JWT-токена
	staticGroup := api.server.Group("/static", api.AuthMiddleware)
	staticGroup.GET("/main.html", func(c echo.Context) error {
		return c.File("../../../static/main.html")
	})
	staticGroup.GET("/style.css", func(c echo.Context) error {
		return c.File("../../../static/style.css")
	})

	api.address = fmt.Sprintf("%s:%d", ip, port)

	return &api
}

func (api *Server) Run() {
	api.server.Logger.Fatal(api.server.Start(api.address))
}
