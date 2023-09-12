package auth

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/logger"
)

type authUC struct {
	cfg      *config.Config
	authRepo postgres.AuthRepository
	log      *logger.ZapLogger
}

func New(cfg *config.Config, authRepo postgres.AuthRepository, log *logger.ZapLogger) *authUC {
	return &authUC{cfg: cfg, authRepo: authRepo, log: log}
}

func (u *authUC) Register() {}
