package auth

import (
	"errors"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/services/chat/internal/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type authRepository struct {
	db *gorm.DB
}

func NewAuthRepository(db *gorm.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(user *models.User) (*models.User, error) {
	// Generate a new UUID for the ID field
	user.ID = uuid.New()
	user.LastLoginAt = time.Now()
	if err := r.db.Create(&user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("email = ?", email).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetByPhone(phone string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("phone = ?", phone).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetByToken(token string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("password_reset_token = ? AND password_reset_at > ?", token, time.Now()).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New(e.ErrTokenExpired)
	}
	return &user, nil
}

func (r *authRepository) GetByCode(code string) (*models.User, error) {
	var user models.User
	if err := r.db.Where("verification_code = ?", code).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, errors.New(e.ErrInvalidToken)
	}
	return &user, nil
}

func (r *authRepository) UpdateLastLogin(userID uuid.UUID) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()")).Error; err != nil {
		return err
	}
	return nil
}
