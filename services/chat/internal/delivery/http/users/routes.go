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
	usersGroup.GET("/account", h.GetAccount())
	usersGroup.GET("/search", h.SearchUsers())
	usersGroup.PUT("/:userId", h.UpdateUser())

	usersGroup.Use(mw.AdminAuthMiddleware())
	usersGroup.PATCH("/:userId/avatar", h.UpdateAvatar())
	usersGroup.DELETE("/:userId", h.DeleteUser())
}
