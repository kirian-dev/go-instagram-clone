package auth

import (
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/pkg/security"
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

func (r *authRepository) GetByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) GetUsers() ([]*models.User, error) {
	var users []*models.User
	if err := r.db.Find(&users).Error; err != nil {
		return nil, err
	}
	for _, user := range users {
		security.DeletePassword(&user.Password)
	}
	return users, nil
}

func (r *authRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *authRepository) DeleteUser(userID uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *authRepository) UpdateLastLogin(userID uuid.UUID) error {
	if err := r.db.Model(&models.User{}).Where("id = ?", userID).Update("last_login_at", gorm.Expr("NOW()")).Error; err != nil {
		return err
	}
	return nil
}
