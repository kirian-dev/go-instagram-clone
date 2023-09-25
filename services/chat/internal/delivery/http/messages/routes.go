package messages

import (
	"go-instagram-clone/services/chat/internal/delivery/http"
	"go-instagram-clone/services/chat/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapMessagesRoutes(messageGroup *echo.Group, h http.MessageHandlers, mw *middleware.MiddlewareManager) {
	messageGroup.Use(mw.AuthJWTMiddleware())
	messageGroup.POST("", h.CreateMessage())
	messageGroup.GET("/list", h.ListMessages())
	messageGroup.GET("/search", h.SearchByText())
	messageGroup.PUT("/:messageID", h.UpdateMessage())
	messageGroup.PATCH("/:messageID", h.ReadMessage())
	messageGroup.DELETE("/:messageID/users/:userID", h.DeleteMessage())
}
