package http

import "github.com/labstack/echo/v4"

type AuthHandlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
	Logout() echo.HandlerFunc
	GetUsers() echo.HandlerFunc
	GetUserByID() echo.HandlerFunc
	DeleteUser() echo.HandlerFunc
	UpdateUser() echo.HandlerFunc
	GetAccount() echo.HandlerFunc
}

type MessageHandlers interface {
	CreateMessage() echo.HandlerFunc
}

type ChatHandlers interface {
	CreateChat() echo.HandlerFunc
	ListChats() echo.HandlerFunc
	DeleteChat() echo.HandlerFunc
	GetChatByID() echo.HandlerFunc
}
