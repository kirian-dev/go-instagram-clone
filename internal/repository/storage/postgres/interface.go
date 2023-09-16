package postgres

import (
	"context"
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)

	GetByEmail(ctx context.Context, email string) (*models.User, error)
	GetByPhone(ctx context.Context, phone string) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
	UpdateUser(ctx context.Context, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type MessagesRepository interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
}
