package chat

import (
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type chatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *chatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) CreateChat(chat *models.Chat) (*models.Chat, error) {
	if err := r.db.Create(chat).Error; err != nil {
		return nil, err
	}
	return chat, nil
}

func (r *chatRepository) ListChats() ([]*models.Chat, error) {
	var chats []*models.Chat
	if err := r.db.Find(&chats).Error; err != nil {
		return nil, err
	}
	return chats, nil
}

func (r *chatRepository) GetChatByID(chatID uuid.UUID) (*models.Chat, error) {
	var chat models.Chat
	if err := r.db.Where("id = ?", chatID).First(&chat).Error; err != nil {
		return nil, err
	}
	return &chat, nil
}

func (r *chatRepository) DeleteChat(chatID uuid.UUID) error {
	if err := r.db.Where("id = ?", chatID).Delete(&models.Chat{}).Error; err != nil {
		return err
	}
	return nil
}
