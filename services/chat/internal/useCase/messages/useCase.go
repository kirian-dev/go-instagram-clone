package messages

import (
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/repository/storage/postgres"
	"time"

	"github.com/google/uuid"
)

type messagesUC struct {
	cfg                 *config.Config
	log                 *logger.ZapLogger
	messagesRepo        postgres.MessagesRepository
	chatParticipantRepo postgres.ChatParticipantRepository
	chatRepo            postgres.ChatRepository
}

func New(cfg *config.Config, messagesRepo postgres.MessagesRepository, chatParticipantRepo postgres.ChatParticipantRepository, chatRepo postgres.ChatRepository, log *logger.ZapLogger) *messagesUC {
	return &messagesUC{cfg, log, messagesRepo, chatParticipantRepo, chatRepo}
}

func (r *messagesUC) CreateMessage(message *models.Message) (*models.Message, error) {
	existsChat, err := r.chatRepo.GetChatByID(message.ChatID)
	if err != nil {
		return nil, err
	}

	if existsChat == nil {
		return nil, errors.New(e.ErrChatNotFound)
	}

	if existsChat.ChatType == models.GroupChat {
		return nil, errors.New(e.ErrNotValidChatType)
	}

	isExistsSender, err := r.chatParticipantRepo.IsParticipantInChat(message.ChatID, message.SenderID)
	if err != nil {
		return nil, err
	}

	if !isExistsSender {
		return nil, errors.New(e.ErrParticipantNotFound)
	}

	isExistsReceiver, err := r.chatParticipantRepo.IsParticipantAdmin(message.ChatID, message.ReceiverID)
	if err != nil {
		return nil, err
	}

	if !isExistsReceiver {
		return nil, errors.New(e.ErrParticipantNotFound)
	}

	createdMessage, err := r.messagesRepo.CreateMessage(message)
	if err != nil {
		return nil, err
	}

	return createdMessage, nil
}

func (r *messagesUC) ListMessages(userID uuid.UUID, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error) {
	messages, err := r.messagesRepo.GetMessages(userID, pag)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messagesUC) SearchByText(userID uuid.UUID, text string, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error) {
	messages, err := r.messagesRepo.SearchByText(userID, text, pag)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (r *messagesUC) UpdateMessage(message *models.Message, messageID, userID uuid.UUID) (*models.Message, error) {
	if message.SenderID != userID {
		return nil, errors.New(e.ErrNoRights)
	}

	existsMessage, err := r.messagesRepo.GetMessageByID(messageID)
	if err != nil {
		return nil, errors.New(e.ErrMessageNotFound)
	}
	existsMessage.Text = message.Text
	existsMessage.UpdatedAt = time.Now()

	updatedMessage, err := r.messagesRepo.UpdateMessage(existsMessage)
	if err != nil {
		return nil, err
	}

	return updatedMessage, nil
}

func (r *messagesUC) ReadMessage(messageID uuid.UUID) (*models.Message, error) {
	existsMessage, err := r.messagesRepo.GetMessageByID(messageID)
	if err != nil {
		return nil, errors.New(e.ErrMessageNotFound)
	}
	existsMessage.ReadAt = true

	updatedMessage, err := r.messagesRepo.UpdateMessage(existsMessage)
	if err != nil {
		return nil, err
	}

	return updatedMessage, nil
}

func (r *messagesUC) DeleteMessage(messageID uuid.UUID) error {
	err := r.messagesRepo.DeleteMessage(messageID)
	if err != nil {
		return err
	}

	return nil
}
