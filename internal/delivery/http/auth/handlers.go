package auth

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/logger"

	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	cfg    *config.Config
	logger logger.ZapLogger
	authUC useCase.AuthUseCase
}

func New(cfg *config.Config, log logger.ZapLogger, authUC useCase.AuthUseCase) *authHandlers {
	return &authHandlers{cfg: cfg, logger: log, authUC: authUC}
}

func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		return nil
	}
}
