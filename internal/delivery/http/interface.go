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
}
