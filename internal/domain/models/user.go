package models

import (
	"go-instagram-clone/pkg/security"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `json:"id" gorm:"primaryKey"`
	FirstName         string    `json:"first_name" validate:"required,lte=50"`
	LastName          string    `json:"last_name" validate:"required,lte=50"`
	Email             string    `json:"email" validate:"required,email,lte=60"`
	Password          string    `json:"password" validate:"required,gte=6"`
	Phone             string    `json:"phone" validate:"omitempty,e164"`
	ProfilePictureURL string    `json:"profile_picture_url" validate:"omitempty,url"`
	City              string    `json:"city" validate:"omitempty,lte=100"`
	Gender            string    `json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday          string    `json:"birthday" validate:"omitempty"`
	Age               int       `json:"age" validate:"omitempty,gte=0,max=200"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Role              string    `json:"role" validate:"omitempty,oneof=user admin"`
	LastLoginAt       time.Time `json:"last_login_at"`
}

func BeforeCreate(u *User) error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	hashedPassword, err := security.HashPassword(u.Password)
	if err != nil {
		return err
	}

	u.Password = hashedPassword

	if u.Phone != "" {
		u.Phone = strings.TrimSpace(u.Phone)
	}
	if u.Role != "" {
		u.Role = strings.ToLower(strings.TrimSpace(u.Role))
	}

	return nil
}
