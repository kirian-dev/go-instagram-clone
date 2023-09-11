package server

import (
	authHttp "go-instagram-clone/internal/delivery/http/auth"
	authRepo "go-instagram-clone/internal/repository/storage/postgres/auth"
	authUseCase "go-instagram-clone/internal/useCase/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Handlers(e *echo.Echo) error {
	// Init repository
	aRepo := authRepo.New(s.db)

	// Init usecase
	authUC := authUseCase.New(s.cfg, aRepo, s.logger)

	// Init delivery
	authHandlers := authHttp.New(s.cfg, s.logger, authUC)

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
	v1 := e.Group("/api/v1")

	authGroup := v1.Group("/auth")
	authHttp.MapAuthRoutes(authGroup, authHandlers)

	return nil
}
