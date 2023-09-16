package useCase

import (
	"context"
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(ctx context.Context, user *models.User) (*models.User, error)
	Login(ctx context.Context, user *models.User) (*models.User, error)
	GetUsers(ctx context.Context) ([]*models.User, error)
	GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User, userID uuid.UUID) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
}

type MessagesUseCase interface {
	CreateMessage(ctx context.Context, message *models.Message) (*models.Message, error)
}
