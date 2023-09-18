package chat

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type chatsHandlers struct {
	cfg    *config.Config
	log    *logger.ZapLogger
	chatUC useCase.ChatsUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, chatUC useCase.ChatsUseCase) *chatsHandlers {
	return &chatsHandlers{cfg, log, chatUC}
}

func (h *chatsHandlers) CreateChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		chat := &models.Chat{}
		if err := utils.ReadRequest(c, chat); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		createdChat, err := h.chatUC.CreateChat(chat)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusCreated, createdChat)
	}
}

func (h *chatsHandlers) ListChats() echo.HandlerFunc {
	return func(c echo.Context) error {
		chats, err := h.chatUC.ListChats()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, chats)
	}
}

func (h *chatsHandlers) GetChatByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDStr := c.Param("chatID")
		chatID, err := uuid.Parse(userIDStr)
		chat, err := h.chatUC.GetChatByID(chatID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, chat)

	}
}

func (h *chatsHandlers) DeleteChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		userIDStr := c.Param("chatID")
		chatID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		if err := h.chatUC.DeleteChat(chatID); err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, map[string]string{"message": "Chat deleted successfully"})
	}
}
