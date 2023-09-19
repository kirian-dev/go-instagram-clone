package useCase

import (
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(user *models.User) (*models.User, error)
	Login(user *models.User) (*models.User, error)
	GetUsers() ([]*models.User, error)
	GetUserByID(userID uuid.UUID) (*models.User, error)
	UpdateUser(user *models.User, userID uuid.UUID) (*models.User, error)
	DeleteUser(userID uuid.UUID) error
}

type MessagesUseCase interface {
	CreateMessage(message *models.Message) (*models.Message, error)
}

type ChatsUseCase interface {
	CreateChatWithParticipants(chatWithParticipants *models.ChatWithParticipants) (*models.ChatWithParticipants, error)
	ListChatsWithParticipants(userID uuid.UUID) ([]*models.ChatWithParticipants, error)
	DeleteChat(chatID uuid.UUID, userID uuid.UUID) error
	GetChatByID(chatID uuid.UUID) (*models.ChatWithParticipants, error)
}
