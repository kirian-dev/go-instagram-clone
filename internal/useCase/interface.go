package useCase

import (
	"context"
	"go-instagram-clone/internal/domain/models"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
}
