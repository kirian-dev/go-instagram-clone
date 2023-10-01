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
	authGroup.POST("/verify-email", h.SendVerificationEmail())
	authGroup.PATCH("/verify-email/:verifyCode", h.VerifyEmail())

	authGroup.Use(mw.AuthJWTMiddleware())
	authGroup.Use(mw.AdminAuthMiddleware())
	authGroup.GET("/analytics/logins", h.GetLoginsCount())
	authGroup.GET("/analytics/registers", h.GetRegistersCount())
}
