package postgres

import (
	"go-instagram-clone/internal/domain/models"

	"github.com/google/uuid"
)

type AuthRepository interface {
	Register(user *models.User) (*models.User, error)

	GetByEmail(email string) (*models.User, error)
	GetByPhone(phone string) (*models.User, error)
	GetUsers() ([]*models.User, error)
	GetByID(userID uuid.UUID) (*models.User, error)
	UpdateLastLogin(userID uuid.UUID) error
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(userID uuid.UUID) error
}

type MessagesRepository interface {
	CreateMessage(message *models.Message) (*models.Message, error)
}

type ChatRepository interface {
	CreateChat(chat *models.Chat) (*models.Chat, error)
	ListChats() ([]*models.Chat, error)
	DeleteChat(chatID uuid.UUID) error
	GetChatByID(chatID uuid.UUID) (*models.Chat, error)
}

type ChatParticipantRepository interface {
	CreateChatParticipant(participant *models.ChatParticipant) error
	GetChatByParticipants([]models.ChatParticipant) (*models.Chat, error)
	GetChatsByUserID(userID uuid.UUID) ([]*models.ChatParticipant, error)
	GetParticipantsByChatID(chatID uuid.UUID) ([]models.ChatParticipant, error)
	DeleteParticipantsByChatID(chatID uuid.UUID) error
	IsParticipantInChat(chatID uuid.UUID, userID uuid.UUID) (bool, error)
}
