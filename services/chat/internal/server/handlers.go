package server

import (
	authHttp "go-instagram-clone/services/chat/internal/delivery/http/auth"
	chatsHttp "go-instagram-clone/services/chat/internal/delivery/http/chats"
	fileImportHttp "go-instagram-clone/services/chat/internal/delivery/http/fileImport"
	messagesHttp "go-instagram-clone/services/chat/internal/delivery/http/messages"
	usersHttp "go-instagram-clone/services/chat/internal/delivery/http/users"
	appMiddleware "go-instagram-clone/services/chat/internal/middleware"
	authRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/auth"
	usersRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/users"
	"go-instagram-clone/services/chat/internal/scheduler"

	chatParticipantsRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/chatParticipants"
	chatsRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/chats"
	fileImportRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/fileImport"

	messagesRepo "go-instagram-clone/services/chat/internal/repository/storage/postgres/messages"
	authUseCase "go-instagram-clone/services/chat/internal/useCase/auth"
	fileImportUseCase "go-instagram-clone/services/chat/internal/useCase/fileImport"
	usersUseCase "go-instagram-clone/services/chat/internal/useCase/users"

	chatsUseCase "go-instagram-clone/services/chat/internal/useCase/chats"
	messagesUseCase "go-instagram-clone/services/chat/internal/useCase/messages"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (s *Server) Handlers(e *echo.Echo) error {
	// Init repository
	aRepo := authRepo.NewAuthRepository(s.db)
	usersRepo := usersRepo.NewUsersRepository(s.db)
	messagesRepo := messagesRepo.NewMessagesRepository(s.db)
	chatRepo := chatsRepo.NewChatRepository(s.db)
	chatParticipantsRepo := chatParticipantsRepo.NewChatParticipantRepository(s.db)
	fileImportRepo := fileImportRepo.NewFileImportRepository(s.db)

	// Init usecase
	authUC := authUseCase.New(s.cfg, aRepo, usersRepo, s.log)
	usersUC := usersUseCase.New(s.cfg, usersRepo, s.log)
	messagesUC := messagesUseCase.New(s.cfg, messagesRepo, chatParticipantsRepo, chatRepo, s.log)
	chatUC := chatsUseCase.New(s.cfg, chatRepo, chatParticipantsRepo, s.log)
	fileImportUC := fileImportUseCase.New(s.cfg, s.log, fileImportRepo, usersRepo, aRepo)

	// Init delivery
	messagesHandlers := messagesHttp.New(s.cfg, s.log, messagesUC)
	authHandlers := authHttp.New(s.cfg, s.log, authUC, s.analyticsClient)
	usersHandlers := usersHttp.New(s.cfg, s.log, usersUC)
	chatsHandlers := chatsHttp.New(s.cfg, s.log, chatUC)
	fileImportHandlers := fileImportHttp.New(s.cfg, s.log, fileImportUC)

	//Api Middleware
	mw := appMiddleware.NewMiddlewareManager(s.cfg, s.log)

	//Swagger
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
	usersGroup := v1.Group("/users")
	usersHttp.MapAuthRoutes(usersGroup, usersHandlers, mw)
	messagesGroup := v1.Group("/messages")
	messagesHttp.MapMessagesRoutes(messagesGroup, messagesHandlers, mw)
	chatsGroup := v1.Group("/chats")
	chatsHttp.MapChatRoutes(chatsGroup, chatsHandlers, mw)
	fileImportGroup := v1.Group("/import")

	fileImportHttp.MapImportRoutes(fileImportGroup, fileImportHandlers, mw)
	scheduler.RunBirthdayCron(s.db, s.log, s.cfg)

	return nil
}
