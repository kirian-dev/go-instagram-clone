package users

import (
	"go-instagram-clone/services/chat/internal/delivery/http"
	"go-instagram-clone/services/chat/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(usersGroup *echo.Group, h http.UsersHandlers, mw *middleware.MiddlewareManager) {
	usersGroup.Use(mw.AuthJWTMiddleware())
	usersGroup.GET("/all", h.GetUsers())
	usersGroup.GET("/:userId", h.GetUserByID())
	usersGroup.DELETE("/:userId", h.DeleteUser())
	usersGroup.PUT("/:userId", h.UpdateUser())
	usersGroup.GET("/account", h.GetAccount())
	usersGroup.PATCH("/:userId/avatar", h.UpdateAvatar())
}
