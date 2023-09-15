package middleware

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
)

// Middleware manager
type MiddlewareManager struct {
	cfg *config.Config
	log *logger.ZapLogger
}

func NewMiddlewareManager(cfg *config.Config, log *logger.ZapLogger) *MiddlewareManager {
	return &MiddlewareManager{cfg: cfg, log: log}
}
