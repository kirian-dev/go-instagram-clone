package messages

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/logger"

	"github.com/labstack/echo/v4"
)

type messagesHandlers struct {
	cfg        *config.Config
	log        *logger.ZapLogger
	messagesUC useCase.MessagesUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, messagesUC useCase.MessagesUseCase) *messagesHandlers {
	return &messagesHandlers{cfg, log, messagesUC}
}

func (h *messagesHandlers) CreateMessage() echo.HandlerFunc {
	return func(c echo.Context) error {
		panic("not implemented")
	}
}
