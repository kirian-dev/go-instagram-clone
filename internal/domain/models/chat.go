package models

import "github.com/google/uuid"

type Chat struct {
	ChatID   uuid.UUID `json:"id" gorm:"primaryKey"`
	ChatName string    `json:"name" validate:"required,max=50"`
}
