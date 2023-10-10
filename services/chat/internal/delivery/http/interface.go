package http

import "github.com/labstack/echo/v4"

type AuthHandlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
	Logout() echo.HandlerFunc
	GetLoginsCount() echo.HandlerFunc
	GetRegistersCount() echo.HandlerFunc
	ForgotPassword() echo.HandlerFunc
	ResetPassword() echo.HandlerFunc
	SendVerificationEmail() echo.HandlerFunc
	VerifyEmail() echo.HandlerFunc
}

type UsersHandlers interface {
	GetUsers() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	GetAccount() echo.HandlerFunc
	UpdateAvatar() echo.HandlerFunc
	SearchUsers() echo.HandlerFunc
}

type MessageHandlers interface {
	CreateMessage() echo.HandlerFunc
	ListMessages() echo.HandlerFunc
	UpdateMessage() echo.HandlerFunc
	ReadMessage() echo.HandlerFunc
	DeleteMessage() echo.HandlerFunc
	SearchByText() echo.HandlerFunc
}

type ChatHandlers interface {
	CreateChatWithParticipants() echo.HandlerFunc
	ListChatsWithParticipants() echo.HandlerFunc
	DeleteChat() echo.HandlerFunc
	GetChatByID() echo.HandlerFunc
	AddParticipantsToChat() echo.HandlerFunc
	RemoveParticipantFromChat() echo.HandlerFunc
}

type FileImportHandlers interface {
	GetImportStatus() echo.HandlerFunc
	GetImportList() echo.HandlerFunc
	UploadFile() echo.HandlerFunc
}
