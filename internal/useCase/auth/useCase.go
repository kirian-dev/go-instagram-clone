package auth

import (
	"context"
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"net/http"
)

type authUC struct {
	cfg      *config.Config
	authRepo postgres.AuthRepository
	log      *logger.ZapLogger
}

func New(cfg *config.Config, authRepo postgres.AuthRepository, log *logger.ZapLogger) *authUC {
	return &authUC{cfg: cfg, authRepo: authRepo, log: log}
}
func (u *authUC) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existsUser, err := u.authRepo.GetByEmail(ctx, user)
	if existsUser != nil || err == nil {
		errorMessage := http.StatusText(http.StatusBadRequest) + ": " + e.ErrInvalidCredentials
		u.log.Error(errorMessage)
		return nil, errors.New(errorMessage)
	}

	if err = user.BeforeCreate(); err != nil {
		u.log.Error("Error in BeforeCreate:", err)
		return nil, err
	}

	newUser, err := u.authRepo.Register(ctx, user)
	if err != nil {
		u.log.Error("Error in Register:", err)
		return nil, err
	}

	utils.DeletePassword(newUser.Password)

	return newUser, nil
}
