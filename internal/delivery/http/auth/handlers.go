package auth

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	cfg       *config.Config
	log       *logger.ZapLogger
	authUC    useCase.AuthUseCase
	validator *validator.Validate
}

func New(cfg *config.Config, log *logger.ZapLogger, authUC useCase.AuthUseCase) *authHandlers {
	validator := validator.New()

	return &authHandlers{cfg: cfg, log: log, authUC: authUC, validator: validator}
}

func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if c.Request().Header.Get("Content-Type") != "application/json" {
			h.log.Error(e.ErrInvalidFormat)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.ErrInvalidFormat})
		}

		user := &models.User{}

		if err := c.Bind(user); err != nil {
			h.log.Error("Error binding user:", err)
			return err
		}
		validationErrors := models.ValidateUser(user, h.validator)
		if len(validationErrors) > 0 {
			h.log.Error("error")
			return c.JSON(http.StatusBadRequest, validationErrors)
		}

		createdUser, err := h.authUC.Register(ctx, user)

		if err != nil {
			h.log.Error(err)
			return err
		}
		return c.JSON(http.StatusCreated, createdUser)
	}
}
