package messages

import (
	"context"
	"database/sql"
	"go-instagram-clone/internal/domain/models"
)

type messagesRepository struct {
	db *sql.DB
}

func NewMessagesRepository(db *sql.DB) *messagesRepository {
	return &messagesRepository{db: db}
}

func (r *messagesRepository) CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error) {
	panic("not implemented")
}
