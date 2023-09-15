package auth

import (
	"go-instagram-clone/internal/delivery/http"
	"go-instagram-clone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h http.AuthHandlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/refresh-token", h.RefreshToken())

	authGroup.Use(mw.AuthJWTMiddleware())
	authGroup.GET("/all", h.GetUsers())
}
