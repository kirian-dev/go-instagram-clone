package postgres

import (
	"context"
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)

	GetByEmail(ctx context.Context, user *models.User) (*models.User, error)
	GetByPhone(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	UpdateLastLogin(ctx context.Context, userID uuid.UUID) error
}
