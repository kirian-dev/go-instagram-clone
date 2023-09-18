package messages

import (
	"go-instagram-clone/internal/domain/models"

	"gorm.io/gorm"
)

type messagesRepository struct {
	db *gorm.DB
}

func NewMessagesRepository(db *gorm.DB) *messagesRepository {
	return &messagesRepository{db: db}
}

func (r *messagesRepository) CreateMessage(message *models.Message) (*models.Message, error) {
	panic("not implemented")
}
