package models

import (
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/utils"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID                string    `db:"id" json:"id"`
	FirstName         string    `db:"first_name" json:"first_name" validate:"required"`
	LastName          string    `db:"last_name" json:"last_name" validate:"required"`
	Email             string    `db:"email" json:"email" validate:"required,email"`
	Password          string    `db:"password" json:"password" validate:"required,min=6"`
	Phone             string    `db:"phone" json:"phone" validate:"omitempty"`
	ProfilePictureURL string    `db:"profile_picture_url" json:"profile_picture_url" validate:"omitempty,url"`
	City              string    `db:"city" json:"city" validate:"omitempty"`
	Gender            string    `db:"gender" json:"gender" validate:"omitempty,oneof=male female other"`
	Birthday          string    `db:"birthday" json:"birthday" validate:"omitempty,date"`
	Age               int       `db:"age" json:"age" validate:"omitempty,gte=0"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	Role              string    `db:"role" json:"role" validate:"omitempty,oneof=user admin"`
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

func ValidateUser(user *User, validate *validator.Validate) []e.ValidationErrorResponse {
	var validationErrors []e.ValidationErrorResponse
	if err := validate.Struct(user); err != nil {
		for _, validationErr := range err.(validator.ValidationErrors) {
			field := validationErr.Field()
			message := ""

			// Customize error message based on the validation tag
			switch validationErr.Tag() {
			case "required":
				message = "Field is required"
			case "email":
				message = "Invalid email format"
			case "min":
				message = "Field length is too short"
			}

			validationErrors = append(validationErrors, e.ValidationErrorResponse{
				Field:   field,
				Message: message,
			})
		}
	}

	return validationErrors
}
