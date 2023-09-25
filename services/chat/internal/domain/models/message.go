package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `json:"id" gorm:"primaryKey"`
	SenderID   uuid.UUID `json:"sender_id" validate:"required,max=50"`
	ReceiverID uuid.UUID `json:"receiver_id" validate:"required,max=50"`
	ChatID     uuid.UUID `json:"chat_id" validate:"required"`
	Text       string    `json:"text" validate:"omitempty,lte=512,required"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
	ReadAt     bool      `json:"is_read"`
}

type MessageListResponse struct {
	Messages   []*Message `json:"messages"`
	TotalCount int64      `json:"totalCount"`
	TotalPages int        `json:"totalPages"`
	Page       int        `json:"page"`
	Size       int        `json:"size"`
	HasMore    bool       `json:"hasMore"`
}
