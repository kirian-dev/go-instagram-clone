package chat

import (
	"go-instagram-clone/internal/delivery/http"
	"go-instagram-clone/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapChatRoutes(chatGroup *echo.Group, h http.ChatHandlers, mw *middleware.MiddlewareManager) {
	chatGroup.Use(mw.AuthJWTMiddleware())
	chatGroup.POST("/", h.CreateChat())
	chatGroup.GET("/:chatID", h.GetChatByID())
	chatGroup.DELETE("/:chatID", h.DeleteChat())
	chatGroup.GET("/list", h.ListChats())
}
