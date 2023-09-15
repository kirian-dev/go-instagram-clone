package models

import (
	"go-instagram-clone/pkg/security"
	"strings"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID                uuid.UUID `db:"id" json:"id"`
	FirstName         string    `db:"first_name" json:"first_name" validate:"required,lte=50"`
	LastName          string    `db:"last_name" json:"last_name" validate:"required,lte=50"`
	Email             string    `db:"email" json:"email" validate:"required,email,lte=60"`
	Password          string    `db:"password" json:"password" validate:"required,gte=6"`
	Phone             string    `db:"phone" json:"phone" validate:"e164,omitempty"`
	ProfilePictureURL string    `db:"profile_picture_url" json:"profile_picture_url" validate:"omitempty,url"`
	City              string    `db:"city" json:"city" validate:"omitempty,lte=100"`
	Gender            string    `db:"gender" json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday          string    `db:"birthday" json:"birthday" validate:"omitempty"`
	Age               int       `db:"age" json:"age" validate:"omitempty,gte=0,max=200"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	Role              string    `db:"role" json:"role" validate:"omitempty,oneof=user admin"`
	LastLoginAt       time.Time `db:"last_login_at" json:"last_login_at"`
}

func (u *User) BeforeCreate() error {
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
