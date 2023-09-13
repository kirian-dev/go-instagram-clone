package postgres

import (
	"context"
	"go-instagram-clone/internal/domain/models"
)

type AuthRepository interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	GetByEmail(ctx context.Context, user *models.User) (*models.User, error)
}
