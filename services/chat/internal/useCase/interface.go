package useCase

import (
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"mime/multipart"

	"github.com/google/uuid"
)

type AuthUseCase interface {
	Register(user *models.User) (*models.UserResponse, error)
	Login(user *models.User) (*models.UserResponse, error)
	ForgotPassword(email string) (*models.UserResponse, string, error)
	ResetPassword(token, password string) error
	SendVerificationEmail(email string) (*models.UserResponse, string, error)
	VerifyEmail(code string) error
}

type UsersUseCase interface {
	GetUsers(pag *utils.PaginationQuery) (*models.UserListResponse, error)
	GetUserByID(userID uuid.UUID) (*models.UserResponse, error)
	UpdateUser(user *models.User, userID uuid.UUID) (*models.UserResponse, error)
	DeleteUser(userID uuid.UUID) error
	UpdateAvatar(userID uuid.UUID, avatarPath string) error
	SearchByQuery(query string, page *utils.PaginationQuery) (*models.UserListResponse, error)
}

type MessagesUseCase interface {
	CreateMessage(message *models.Message) (*models.Message, error)
	ListMessages(userID uuid.UUID, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error)
	ReadMessage(messageID uuid.UUID) (*models.Message, error)
	DeleteMessage(messageID uuid.UUID) error
	UpdateMessage(message *models.Message, messageID, userID uuid.UUID) (*models.Message, error)
	SearchByText(userID uuid.UUID, text string, pag *utils.PaginationQuery) ([]*models.MessageListResponse, error)
}

type ChatsUseCase interface {
	CreateChatWithParticipants(chatWithParticipants *models.ChatWithParticipants) (*models.ChatWithParticipants, error)
	ListChatsWithParticipants(userID uuid.UUID) ([]*models.ChatWithParticipants, error)
	DeleteChat(chatID uuid.UUID, userID uuid.UUID) error
	GetChatByID(chatID uuid.UUID) (*models.ChatWithParticipants, error)
	AddParticipantsToChat(participants []*models.ChatParticipant, chatID uuid.UUID, userID uuid.UUID) ([]*models.ChatParticipant, error)
	RemoveParticipantFromChat(chatID, userID, participantID uuid.UUID) error
}

type FileImportUseCase interface {
	UploadFile(file *multipart.FileHeader) error
	GetImportFiles() ([]*models.FileImport, error)
	GetImportFileByID(fileID uuid.UUID) (*models.FileImport, error)
}
