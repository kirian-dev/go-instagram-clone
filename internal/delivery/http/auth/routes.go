package auth

import (
	"go-instagram-clone/internal/delivery/http"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h http.AuthHandlers) {
	authGroup.POST("/register", h.Register())
}
