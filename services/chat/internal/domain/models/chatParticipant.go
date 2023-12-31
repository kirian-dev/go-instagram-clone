package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatParticipant struct {
	ChatParticipantID uuid.UUID       `json:"id" gorm:"primaryKey"`
	ChatID            uuid.UUID       `json:"chat_id" validate:"required"`
	UserID            uuid.UUID       `json:"user_id" validate:"required"`
	Role              ParticipantRole `json:"role" validate:"required"`
	JoinedAt          time.Time       `json:"joined_at"`
}

type ParticipantRole string

const (
	Admin  ParticipantRole = "admin"
	Member ParticipantRole = "member"
)

func (cp *ChatParticipant) BeforeCreate(tx *gorm.DB) (err error) {
	cp.ChatParticipantID = uuid.New()
	cp.JoinedAt = time.Now()
	return nil
}
