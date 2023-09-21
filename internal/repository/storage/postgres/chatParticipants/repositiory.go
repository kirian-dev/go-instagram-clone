package chatParticipants

import (
	"errors"
	"go-instagram-clone/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatParticipantRepository struct {
	db *gorm.DB
}

func NewChatParticipantRepository(db *gorm.DB) *ChatParticipantRepository {
	return &ChatParticipantRepository{db}
}

func (r *ChatParticipantRepository) CreateChatParticipant(participant *models.ChatParticipant) error {
	participant.ChatParticipantID = uuid.New()
	participant.JoinedAt = time.Now()
	if err := r.db.Create(participant).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChatParticipantRepository) GetChatByParticipants(participants []models.ChatParticipant) (*models.Chat, error) {
	var userIds []uuid.UUID
	for _, participant := range participants {
		userIds = append(userIds, participant.UserID)
	}

	// Query the database to find a chat with the given participants
	var chat models.Chat
	result := r.db.Table("chat_participants").
		Select("chat_id").
		Where("user_id IN (?)", userIds).
		Group("chat_id").
		Having("COUNT(DISTINCT user_id) = ?", len(participants)).
		First(&chat)

	if result.Error != nil && result.RowsAffected == 0 {
		return nil, result.Error
	}

	// Load the chat details
	if err := r.db.Where("chat_id = ?", chat.ChatID).First(&chat).Error; err != nil {
		return nil, err
	}

	return &chat, nil
}

func (r *ChatParticipantRepository) GetParticipantsByChatID(chatID uuid.UUID) ([]models.ChatParticipant, error) {
	var participants []models.ChatParticipant
	if err := r.db.Where("chat_id = ?", chatID).Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ChatParticipantRepository) GetChatsByUserID(userID uuid.UUID) ([]*models.ChatParticipant, error) {
	var participants []*models.ChatParticipant
	if err := r.db.Where("user_id = ?", userID).Find(&participants).Error; err != nil {
		return nil, err
	}
	return participants, nil
}

func (r *ChatParticipantRepository) DeleteParticipantsByChatID(chatID uuid.UUID) error {
	if err := r.db.Where("chat_id = ?", chatID).Delete(&models.ChatParticipant{}).Error; err != nil {
		return err
	}
	return nil
}

func (r *ChatParticipantRepository) IsParticipantInChat(chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	var participant models.ChatParticipant
	result := r.db.Where("chat_id = ? AND user_id = ?", chatID, userID).First(&participant)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (r *ChatParticipantRepository) IsParticipantAdmin(chatID uuid.UUID, userID uuid.UUID) (bool, error) {
	var participant models.ChatParticipant
	result := r.db.Where("chat_id = ? AND user_id = ? AND role = ?", chatID, userID, models.Admin).First(&participant)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return false, nil
	}

	return true, nil
}

func (r *ChatParticipantRepository) DeleteParticipantFromChat(chatID, participantID uuid.UUID) error {
	if err := r.db.Where("chat_id = ? AND user_id = ?", chatID, participantID).Delete(&models.ChatParticipant{}).Error; err != nil {
		return err
	}
	return nil
}
