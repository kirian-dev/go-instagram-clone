package models

import (
	"time"

	"github.com/google/uuid"
)

type Chat struct {
	ChatID    uuid.UUID `json:"id" gorm:"primaryKey"`
	ChatName  string    `json:"name" validate:"omitempty,max=50,unique"`
	ChatType  ChatType  `json:"type" validate:"required"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type ChatType string

const (
	PrivateChat ChatType = "private"
	GroupChat   ChatType = "group"
)

type ChatWithParticipants struct {
	Chat         Chat              `json:"chat"`
	Participants []ChatParticipant `json:"participants"`
}
