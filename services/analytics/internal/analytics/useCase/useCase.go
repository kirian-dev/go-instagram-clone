package useCase

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/services/analytics/internal/analytics/repository/storage/postgres"
)

type AnalyticsUseCase interface {
	RecordLogin(email, phone string) error
	RecordNewUser(email, phone string) error
	GetQuantityLogins() (int32, error)
	GetQuantityRegister() (int32, error)
}

type analyticsUC struct {
	cfg   *config.Config
	log   *logger.ZapLogger
	aRepo postgres.AnalyticsRepository
}

func NewAnalyticsUC(cfg *config.Config, log *logger.ZapLogger, aRepo postgres.AnalyticsRepository) *analyticsUC {
	return &analyticsUC{cfg: cfg, log: log, aRepo: aRepo}
}

func (uc *analyticsUC) RecordLogin(email, phone string) error {
	err := uc.aRepo.SaveSuccessfulLogin(email, phone)
	if err != nil {
		return err
	}

	return nil
}

func (uc *analyticsUC) RecordNewUser(email, phone string) error {
	err := uc.aRepo.SaveSuccessfulRegister(email, phone)
	if err != nil {
		return err
	}

	return nil
}

func (uc *analyticsUC) GetQuantityLogins() (int32, error) {
	quantity, err := uc.aRepo.GetQuantityLogins()
	if err != nil {
		return 0, err
	}

	return quantity, nil
}

func (uc *analyticsUC) GetQuantityRegister() (int32, error) {
	quantity, err := uc.aRepo.GetQuantityRegister()
	if err != nil {
		return 0, err
	}

	return quantity, nil
}
