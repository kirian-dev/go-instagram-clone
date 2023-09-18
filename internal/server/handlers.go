package server

import (
	authHttp "go-instagram-clone/internal/delivery/http/auth"
	chatsHttp "go-instagram-clone/internal/delivery/http/chats"
	messagesHttp "go-instagram-clone/internal/delivery/http/messages"

	appMiddleware "go-instagram-clone/internal/middleware"
	authRepo "go-instagram-clone/internal/repository/storage/postgres/auth"
	chatsRepo "go-instagram-clone/internal/repository/storage/postgres/chats"
	messagesRepo "go-instagram-clone/internal/repository/storage/postgres/messages"
	authUseCase "go-instagram-clone/internal/useCase/auth"
	chatsUseCase "go-instagram-clone/internal/useCase/chats"
	messagesUseCase "go-instagram-clone/internal/useCase/messages"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) Handlers(e *echo.Echo) error {
	// Init repository
	aRepo := authRepo.NewAuthRepository(s.db)
	messagesRepo := messagesRepo.NewMessagesRepository(s.db)
	chatRepo := chatsRepo.NewChatRepository(s.db)

	// Init usecase
	authUC := authUseCase.New(s.cfg, aRepo, s.log)
	messagesUC := messagesUseCase.New(s.cfg, messagesRepo, s.log)
	chatUC := chatsUseCase.New(s.cfg, chatRepo, s.log)

	// Init delivery
	messagesHandlers := messagesHttp.New(s.cfg, s.log, messagesUC)
	authHandlers := authHttp.New(s.cfg, s.log, authUC)
	chatsHandlers := chatsHttp.New(s.cfg, s.log, chatUC)

	//Api Middleware
	mw := appMiddleware.NewMiddlewareManager(s.cfg, s.log)

	e.GET("/swagger/*", echoSwagger.WrapHandler)

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

	chatsGroup := v1.Group("/chats")
	chatsHttp.MapChatRoutes(chatsGroup, chatsHandlers, mw)
	//Swagger

	return nil
}
