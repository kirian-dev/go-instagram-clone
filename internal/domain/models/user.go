package models

import (
	"go-instagram-clone/pkg/utils"
	"strings"
	"time"
)

type User struct {
	ID                string    `db:"id" json:"id"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	Email             string    `db:"email" json:"email"`
	Password          string    `db:"password" json:"password"`
	Phone             string    `db:"phone" json:"phone"`
	ProfilePictureURL string    `db:"profile_picture_url" json:"profile_picture_url"`
	City              string    `db:"city" json:"city"`
	Gender            string    `db:"gender" json:"gender"`
	Birthday          string    `db:"birthday" json:"birthday"`
	Age               int       `db:"age" json:"age"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	Role              string    `db:"role" json:"role"`
	LastLoginAt       time.Time `db:"last_login_at" json:"last_login_at"`
}

func (u *User) BeforeCreate() error {
	u.Email = strings.ToLower(strings.TrimSpace(u.Email))
	u.Password = strings.TrimSpace(u.Password)

	if err := utils.HashPassword(u.Password); err != nil {
		return err
	}

	if u.Phone != "" {
		u.Phone = strings.TrimSpace(u.Phone)
	}
	if u.Role != "" {
		u.Role = strings.ToLower(strings.TrimSpace(u.Role))
	}

	return nil
}
