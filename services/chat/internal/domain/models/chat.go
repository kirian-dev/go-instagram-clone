package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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

func (c *Chat) BeforeCreate(tx *gorm.DB) (err error) {
	c.ChatID = uuid.New()
	c.CreatedAt = time.Now()
	c.UpdatedAt = time.Now()
	return nil
}

func (c *Chat) BeforeUpdate(tx *gorm.DB) (err error) {
	c.UpdatedAt = time.Now()
	return nil
}
