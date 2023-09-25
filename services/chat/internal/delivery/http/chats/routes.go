package chat

import (
	"go-instagram-clone/services/chat/internal/delivery/http"
	"go-instagram-clone/services/chat/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapChatRoutes(chatGroup *echo.Group, h http.ChatHandlers, mw *middleware.MiddlewareManager) {
	chatGroup.Use(mw.AuthJWTMiddleware())
	chatGroup.POST("", h.CreateChatWithParticipants())
	chatGroup.GET("/:chatID", h.GetChatByID())
	chatGroup.DELETE("/:chatID", h.DeleteChat())
	chatGroup.GET("/list", h.ListChatsWithParticipants())
	chatGroup.POST("/:chatID/participants", h.AddParticipantsToChat())
	chatGroup.DELETE("/:chatID/participants/:participantID", h.RemoveParticipantFromChat())
}
