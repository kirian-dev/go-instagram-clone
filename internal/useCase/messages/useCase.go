package messages

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/logger"
)

type messagesUC struct {
	cfg          *config.Config
	log          *logger.ZapLogger
	messagesRepo postgres.MessagesRepository
}

func New(cfg *config.Config, messagesRepo postgres.MessagesRepository, log *logger.ZapLogger) *messagesUC {
	return &messagesUC{cfg, log, messagesRepo}
}

func (r *messagesUC) CreateMessage(message *models.Message) (*models.Message, error) {
	panic("not implemented")
}
