package messages

import (
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type messagesRepository struct {
	db *gorm.DB
}

func NewMessagesRepository(db *gorm.DB) *messagesRepository {
	return &messagesRepository{db: db}
}

func (r *messagesRepository) CreateMessage(message *models.Message) (*models.Message, error) {
	message.ReadAt = false
	message.ID = uuid.New()

	if err := r.db.Create(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *messagesRepository) UpdateMessage(message *models.Message) (*models.Message, error) {
	if err := r.db.Save(message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *messagesRepository) GetMessages(userID uuid.UUID) ([]*models.Message, error) {
	var messages []*models.Message

	if err := r.db.Where("sender_id = ?", userID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *messagesRepository) GetMessageByID(messageID uuid.UUID) (*models.Message, error) {
	var message *models.Message

	if err := r.db.Where("id = ?", messageID).First(&message).Error; err != nil {
		return nil, err
	}
	return message, nil
}

func (r *messagesRepository) DeleteMessage(messageID uuid.UUID) error {
	if err := r.db.Delete(&models.Message{}, messageID).Error; err != nil {
		return err
	}
	return nil
}
