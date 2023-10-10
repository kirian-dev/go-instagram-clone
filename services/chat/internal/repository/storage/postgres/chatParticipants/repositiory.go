package chatParticipants

import (
	"go-instagram-clone/services/chat/internal/domain/models"

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

	var chat models.Chat
	query := r.db.Where("user_id IN (?)", userIds).Group("chat_id").Having("COUNT(DISTINCT user_id) = ?", len(participants))
	result := query.Table("chat_participants").Select("chat_id").First(&chat)

	if result.Error != nil {
		return nil, result.Error
	}

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

func (r *ChatParticipantRepository) IsParticipantInChat(chatID, userID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Table("chat_participants").
		Where("chat_id = ? AND user_id = ?", chatID, userID).
		Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

func (r *ChatParticipantRepository) IsParticipantAdmin(chatID, userID uuid.UUID) (bool, error) {
	var count int64
	query := r.db.Table("chat_participants").
		Where("chat_id = ? AND user_id = ? AND role = ?", chatID, userID, models.Admin).
		Count(&count)
	if query.Error != nil {
		return false, query.Error
	}
	return count > 0, nil
}

func (r *ChatParticipantRepository) DeleteParticipantFromChat(chatID, participantID uuid.UUID) error {
	if err := r.db.Where("chat_id = ? AND user_id = ?", chatID, participantID).Delete(&models.ChatParticipant{}).Error; err != nil {
		return err
	}
	return nil
}
