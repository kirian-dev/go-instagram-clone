package models

import (
	"time"

	"github.com/google/uuid"
)

type Message struct {
	ID         uuid.UUID `json:"id" db:"id"`
	SenderID   uuid.UUID `json:"sender_id" db:"sender_id" validate:"required,max=50"`
	ReceiverID uuid.UUID `json:"receiver_id" db:"receiver_id" validate:"required,max=50"`
	Text       string    `json:"text" db:"text" validate:"omitempty,lte=512,required"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	ReadAt     time.Time `json:"read_at" db:"read_at" `
}
