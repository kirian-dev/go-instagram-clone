package chat

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/useCase"
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

// @Summary Create chat
// @Description Create a private or group chat
// @Accept json
// @Produce json
// @Tags Chats
// @Success 201 {object} models.ChatWithParticipants
// @Failure 400 {object}  e.ErrorResponse "Error creating chat"
func (h *chatsHandlers) CreateChatWithParticipants() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)
		chatWithParticipants := &models.ChatWithParticipants{}
		if err := utils.ReadRequest(c, chatWithParticipants); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		roleMap := make(map[uuid.UUID]models.ParticipantRole)
		for _, participant := range chatWithParticipants.Participants {
			roleMap[participant.UserID] = participant.Role
		}

		currentUserRole, hasAccess := roleMap[userClaims.UserID]
		if !hasAccess || currentUserRole != models.Admin {
			return c.JSON(http.StatusForbidden, e.ErrorResponse{Error: e.ErrCreateChatNoRights})
		}

		// Create the chat with participants.
		createdChat, err := h.chatUC.CreateChatWithParticipants(chatWithParticipants)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusCreated, createdChat)
	}
}

// @Summary Get all chats
// @Description List chats for a current user
// @Tags Chats
// @Accept  json
// @Produce  json
// @Success 200 {object} models.ChatWithParticipants
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /chats/list [get]
func (h *chatsHandlers) ListChatsWithParticipants() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)

		chats, err := h.chatUC.ListChatsWithParticipants(userClaims.UserID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, chats)
	}
}

// @Summary Get chat
// @Get chat
// @Tags Chats
// @Accept  json
// @Produce  json
// @Param chatID path int true "chatID"
// @Success 200 {object} models.ChatWithParticipants
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /chats/{chatID} [get]
func (h *chatsHandlers) GetChatByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		chatIDStr := c.Param("chatID")
		chatID, err := uuid.Parse(chatIDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, err)
		}
		chat, err := h.chatUC.GetChatByID(chatID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusOK, chat)

	}
}

// @Summary Delete chat
// @Description Delete chat if user is admin or member in this chat
// @Tags Chats
// @Param chatID path int true "chatID"
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /chats/{chatID} [delete]
func (h *chatsHandlers) DeleteChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)
		chatDStr := c.Param("chatID")
		chatID, err := uuid.Parse(chatDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		if err := h.chatUC.DeleteChat(chatID, userClaims.UserID); err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusNoContent, map[string]string{"message": "Chat deleted successfully"})
	}
}

// @Summary Add Participants to group chat
// @Description Add participants to group chat
// @Tags Chats
// @Param chatID path int true "chatID"
// @Produce json
// @Success 200 {object} models.ChatParticipant
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /{chatID}/participants [post]
func (h *chatsHandlers) AddParticipantsToChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)
		chatDStr := c.Param("chatID")
		chatID, err := uuid.Parse(chatDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		participants := []*models.ChatParticipant{}

		if err := utils.ReadRequest(c, &participants); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		addedParticipants, err := h.chatUC.AddParticipantsToChat(participants, chatID, userClaims.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, addedParticipants)
	}
}

// @Summary Delete participant from chat
// @Description Delete participant from group chat
// @Tags Chats
// @Param chatID path int true "chatID"
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /chats/{chatID}/participants/{participantID} [delete]
func (h *chatsHandlers) RemoveParticipantFromChat() echo.HandlerFunc {
	return func(c echo.Context) error {
		chatDStr := c.Param("chatID")
		participantDStr := c.Param("participantID")

		chatID, err := uuid.Parse(chatDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		participantID, err := uuid.Parse(participantDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		userClaims := c.Get("userClaims").(*utils.CustomClaims)

		if err := h.chatUC.RemoveParticipantFromChat(chatID, userClaims.UserID, participantID); err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		return c.JSON(http.StatusNoContent, map[string]string{"message": "Participant deleted successfully"})
	}
}
