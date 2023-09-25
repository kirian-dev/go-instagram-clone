package messages

import (
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/pkg/utils"

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

func (r *messagesRepository) GetMessages(userID uuid.UUID, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error) {
	offset := pag.GetOffset()
	limit := pag.GetLimit()

	var totalCount int64

	if err := r.db.Model(&models.Message{}).
		Where("sender_id = ?", userID).
		Count(&totalCount).
		Error; err != nil {
		return nil, err
	}

	var messages []*models.Message

	query := r.db.Where("sender_id = ?", userID).Offset(offset).Limit(limit)
	if err := query.Find(&messages).Error; err != nil {
		return nil, err
	}

	response := []*models.MessageListResponse{
		{
			Messages:   messages,
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pag.GetSize()),
			Page:       pag.GetPage(),
			Size:       pag.GetSize(),
			HasMore:    utils.GetHasMore(pag.GetPage(), totalCount, pag.GetSize()),
		},
	}

	return response, nil
}

func (r *messagesRepository) SearchByText(userID uuid.UUID, text string, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error) {
	offset := pag.GetOffset()
	limit := pag.GetLimit()

	var totalCount int64

	// Count total messages that match the query
	if err := r.db.Model(&models.Message{}).
		Where("sender_id = ? AND text LIKE ?", userID, "%"+text+"%").
		Count(&totalCount).
		Error; err != nil {
		return nil, err
	}

	var messages []*models.Message

	if err := r.db.Model(&models.Message{}).
		Select("sender_id, text").
		Where("sender_id = ? AND text LIKE ?", userID, "%"+text+"%").
		Offset(offset).
		Limit(limit).
		Find(&messages).
		Error; err != nil {
		return nil, err
	}

	response := []*models.MessageListResponse{
		{
			Messages:   messages,
			TotalCount: totalCount,
			TotalPages: utils.GetTotalPages(totalCount, pag.GetSize()),
			Page:       pag.GetPage(),
			Size:       pag.GetSize(),
			HasMore:    utils.GetHasMore(pag.GetPage(), totalCount, pag.GetSize()),
		},
	}

	return response, nil
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
