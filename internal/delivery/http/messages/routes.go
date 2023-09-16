package messages

import (
	"go-instagram-clone/internal/delivery/http"
	"go-instagram-clone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapMessagesRoutes(messageGroup *echo.Group, h http.MessageHandlers, mw *middleware.MiddlewareManager) {
	messageGroup.Use(mw.AuthJWTMiddleware())
	messageGroup.POST("/", h.CreateMessage())
}
