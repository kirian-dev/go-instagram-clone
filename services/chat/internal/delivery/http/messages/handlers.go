package messages

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

type messagesHandlers struct {
	cfg        *config.Config
	log        *logger.ZapLogger
	messagesUC useCase.MessagesUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, messagesUC useCase.MessagesUseCase) *messagesHandlers {
	return &messagesHandlers{cfg, log, messagesUC}
}

// @Summary Create message
// @Description Create message for private chat
// @Accept json
// @Produce json
// @Tags Messages
// @Success 201 {object} models.Message
// @Failure 400 {object} e.ErrorResponse "Error creating message"
func (h *messagesHandlers) CreateMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		message := &models.Message{}
		userClaims := c.Get("userClaims").(*utils.CustomClaims)

		if err := utils.ReadRequest(c, message); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		if message.SenderID != userClaims.UserID {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: e.ErrNotCorrectSender})
		}

		createdMessage, err := h.messagesUC.CreateMessage(message)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusCreated, createdMessage)
	}
}

// @Summary Get all messages
// @Description List messages for a current user with pagination
// @Tags Messages
// @Accept  json
// @Produce  json
// @Success 200 {object} models.MessageListResponse
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /messages/list [get]
func (h *messagesHandlers) ListMessages() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)

		pag, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		messages, err := h.messagesUC.ListMessages(userClaims.UserID, pag)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, messages)
	}
}

// @Summary Search for messages
// @Description  search for messages
// @Tags Messages
// @Accept  json
// @Produce  json
// @Success 200 {object} models.MessageListResponse
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /messages/search [get]
func (h *messagesHandlers) SearchByText() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)

		pag, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		messages, err := h.messagesUC.SearchByText(userClaims.UserID, c.QueryParam("text"), pag)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, messages)
	}
}

// @Summary Update message
// @Description Update message only for who created the message
// @Tags Messages
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Message
// @Param messageID path int true "messageID"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /messages/{messageID} [put]
func (h *messagesHandlers) UpdateMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)
		messageDStr := c.Param("messageID")
		messageID, err := uuid.Parse(messageDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		message := &models.Message{}
		if err := utils.ReadRequest(c, message); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		updatedMessage, err := h.messagesUC.UpdateMessage(message, messageID, userClaims.UserID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedMessage)
	}
}

// @Summary Update read message
// @Description Update read at message
// @Tags Messages
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Message
// @Param messageID path int true "messageID"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /messages/{messageID} [patch]
func (h *messagesHandlers) ReadMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		messageDStr := c.Param("messageID")
		messageID, err := uuid.Parse(messageDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		updatedMessage, err := h.messagesUC.ReadMessage(messageID)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedMessage)
	}
}

// @Summary Delete message
// @Description Delete message if user is created this message
// @Tags Messages
// @Param chatID path int true "chatID"
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /messages/{messageID}/users/{userID} [delete]
func (h *messagesHandlers) DeleteMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims := c.Get("userClaims").(*utils.CustomClaims)
		messageDStr := c.Param("messageID")
		messageID, err := uuid.Parse(messageDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		userIDStr := c.Param("userID")
		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		if userID != userClaims.UserID {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: e.ErrNoRights})
		}

		err = h.messagesUC.DeleteMessage(messageID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusNoContent, map[string]string{"message": "message successfully deleted"})
	}
}
