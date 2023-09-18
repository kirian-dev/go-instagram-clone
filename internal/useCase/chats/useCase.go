package chat

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/logger"

	"github.com/google/uuid"
)

type chatsUC struct {
	cfg      *config.Config
	log      *logger.ZapLogger
	chatRepo postgres.ChatRepository
}

func New(cfg *config.Config, chatRepo postgres.ChatRepository, log *logger.ZapLogger) *chatsUC {
	return &chatsUC{cfg, log, chatRepo}
}

func (uc *chatsUC) CreateChat(chat *models.Chat) (*models.Chat, error) {
	createdChat, err := uc.chatRepo.CreateChat(chat)
	if err != nil {
		return nil, err
	}
	return createdChat, nil
}

func (uc *chatsUC) ListChats() ([]*models.Chat, error) {
	chats, err := uc.chatRepo.ListChats()
	if err != nil {
		return nil, err
	}
	return chats, nil
}

func (uc *chatsUC) GetChatByID(chatID uuid.UUID) (*models.Chat, error) {
	chat, err := uc.chatRepo.GetChatByID(chatID)
	if err != nil {
		return nil, err
	}
	return chat, nil
}

func (uc *chatsUC) DeleteChat(chatID uuid.UUID) error {
	err := uc.chatRepo.DeleteChat(chatID)
	if err != nil {
		return err
	}
	return nil
}
