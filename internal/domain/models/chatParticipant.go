package models

import "github.com/google/uuid"

type ChatParticipant struct {
	ChatParticipantID uuid.UUID `json:"id" gorm:"primaryKey"`
	ChatID            uuid.UUID `json:"chat_id" validate:"required" gorm:"foreignKey"`
	UserID            uuid.UUID `json:"user_id" validate:"required" gorm:"foreignKey"`
}
