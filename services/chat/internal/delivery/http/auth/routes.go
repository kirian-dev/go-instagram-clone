package auth

import (
	"go-instagram-clone/services/chat/internal/delivery/http"
	"go-instagram-clone/services/chat/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapAuthRoutes(authGroup *echo.Group, h http.AuthHandlers, mw *middleware.MiddlewareManager) {
	authGroup.POST("/register", h.Register())
	authGroup.POST("/login", h.Login())
	authGroup.POST("/refresh-token", h.RefreshToken())
	authGroup.POST("/logout", h.Logout())
	authGroup.POST("/forgot-password", h.ForgotPassword())
	authGroup.PATCH("/reset-password/:resetToken", h.ResetPassword())

	authGroup.Use(mw.AuthJWTMiddleware())
	authGroup.GET("/all", h.GetUsers())
	authGroup.GET("/:userId", h.GetUserByID())
	authGroup.DELETE("/:userId", h.DeleteUser())
	authGroup.PUT("/:userId", h.UpdateUser())
	authGroup.GET("/account", h.GetAccount())

	authGroup.Use(mw.AdminAuthMiddleware())
	authGroup.GET("/analytics/logins", h.GetLoginsCount())
	authGroup.GET("/analytics/registers", h.GetRegistersCount())
}
