package auth

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type authHandlers struct {
	cfg    *config.Config
	log    *logger.ZapLogger
	authUC useCase.AuthUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, authUC useCase.AuthUseCase) *authHandlers {
	return &authHandlers{cfg: cfg, log: log, authUC: authUC}
}

func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		user := &models.User{}
		createdUser, err := h.authUC.Register(ctx, user)

		if err != nil {
			log.Error(err)
			return err
		}
		return c.JSON(http.StatusCreated, createdUser)
	}
}
