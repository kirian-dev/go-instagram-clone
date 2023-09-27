package useCase

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/services/analytics/internal/analytics/repository/storage/mysql"
)

type AnalyticsUseCase interface {
	RecordLogin(logins int32) error
	RecordNewUser(registers int32) error
	GetQuantityLogins() (int32, error)
	GetQuantityRegister() (int32, error)
}

type analyticsUC struct {
	cfg   *config.Config
	log   *logger.ZapLogger
	aRepo mysql.AnalyticsRepository
}

func NewAnalyticsUC(cfg *config.Config, log *logger.ZapLogger, aRepo mysql.AnalyticsRepository) *analyticsUC {
	return &analyticsUC{cfg: cfg, log: log, aRepo: aRepo}
}

func (uc *analyticsUC) RecordLogin(logins int32) error {
	err := uc.aRepo.SaveSuccessfulLogin(logins)
	if err != nil {
		return err
	}

	return nil
}

func (uc *analyticsUC) RecordNewUser(registers int32) error {
	err := uc.aRepo.SaveSuccessfulRegister(registers)
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
