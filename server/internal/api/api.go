package api

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var jwtSecret = []byte("your-secret-key") // Замените на свой секретный ключ

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

	// Настройка CORS для всего сервера
	api.server.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},                                                                                    // Разрешаем запросы со всех источников (для разработки)
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions},                    // Добавляем PUT
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization}, // Добавляем Authorization
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
	staticGroup := api.server.Group("/static", AuthMiddleware)
	staticGroup.GET("/main.html", func(c echo.Context) error {
		return c.File("../../../static/main.html")
	})
	staticGroup.GET("/style.css", func(c echo.Context) error {
		return c.File("../../../static/style.css")
	})
	// Добавьте другие статические файлы, если необходимо

	api.address = fmt.Sprintf("%s:%d", ip, port)

	return &api
}

func (api *Server) Run() {
	api.server.Logger.Fatal(api.server.Start(api.address))
}

// Черный список для отозванных токенов
var revokedTokens = struct {
	sync.RWMutex
	tokens map[string]bool
}{tokens: make(map[string]bool)}

// Middleware для проверки токена
func AuthMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		tokenString := strings.Replace(authHeader, "Bearer ", "", 1)

		// Проверяем, находится ли токен в черном списке
		revokedTokens.RLock()
		if revokedTokens.tokens[tokenString] {
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
		if err != nil {
			return c.Redirect(http.StatusTemporaryRedirect, "/login")
		}

		if claims, ok := token.Claims.(*jwtCustomClaims); ok && token.Valid {
			c.Set("email", claims.Email)
			return next(c)
		}

		return c.Redirect(http.StatusTemporaryRedirect, "/login")
	}
}
