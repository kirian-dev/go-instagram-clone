package http

import "github.com/labstack/echo/v4"

type AuthHandlers interface {
	Register() echo.HandlerFunc
	GetUsers() echo.HandlerFunc
	Login() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
}
