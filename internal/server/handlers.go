package server

import (
	authHttp "go-instagram-clone/internal/delivery/http/auth"
	messagesHttp "go-instagram-clone/internal/delivery/http/messages"
	appMiddleware "go-instagram-clone/internal/middleware"
	authRepo "go-instagram-clone/internal/repository/storage/postgres/auth"
	authUseCase "go-instagram-clone/internal/useCase/auth"

	messagesRepo "go-instagram-clone/internal/repository/storage/postgres/messages"
	messagesUseCase "go-instagram-clone/internal/useCase/messages"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) Handlers(e *echo.Echo) error {
	// Init repository
	aRepo := authRepo.NewAuthRepository(s.db)
	messagesRepo := messagesRepo.NewMessagesRepository(s.db)

	// Init usecase
	authUC := authUseCase.New(s.cfg, aRepo, s.log)
	messagesUC := messagesUseCase.New(s.cfg, messagesRepo, s.log)

	// Init delivery
	messagesHandlers := messagesHttp.New(s.cfg, s.log, messagesUC)
	authHandlers := authHttp.New(s.cfg, s.log, authUC)

	//Api Middleware
	mw := appMiddleware.NewMiddlewareManager(s.cfg, s.log)

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
	authHttp.MapAuthRoutes(authGroup, authHandlers, mw)

	messagesGroup := v1.Group("/messages")
	messagesHttp.MapMessagesRoutes(messagesGroup, messagesHandlers, mw)
	return nil
}
