package mysql

import (
	"go-instagram-clone/services/analytics/internal/models"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	SaveSuccessfulLogin(logins int32) error
	SaveSuccessfulRegister(registers int32) error
	GetQuantityLogins() (int32, error)
	GetQuantityRegister() (int32, error)
}

type analyticsRepo struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *analyticsRepo {
	return &analyticsRepo{db: db}
}

func (r *analyticsRepo) SaveSuccessfulLogin(logins int32) error {
	// Start a transaction to ensure atomicity
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Retrieve the current total logins
	var currentTotalLogins int32
	if err := tx.Model(&models.Analytics{}).Select("successful_logins").Scan(&currentTotalLogins).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the total logins
	newTotalLogins := currentTotalLogins + logins
	if err := tx.Model(&models.Analytics{}).Update("successful_logins", newTotalLogins).Error; err != nil {
		tx.Rollback()
		return err
	}

	analytics := &models.Analytics{
		SuccessfulLogins:            logins,
		SuccessfulLoginLastUpdateAt: time.Now(),
	}
	if err := tx.Create(analytics).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *analyticsRepo) SaveSuccessfulRegister(registers int32) error {
	// Start a transaction to ensure atomicity
	tx := r.db.Begin()
	if tx.Error != nil {
		return tx.Error
	}

	// Retrieve the current total register
	var currentTotalRegisters int32
	if err := tx.Model(&models.Analytics{}).Select("successful_register").Scan(&currentTotalRegisters).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Update the total registers
	newTotalRegisters := currentTotalRegisters + registers
	if err := tx.Model(&models.Analytics{}).Update("successful_register", newTotalRegisters).Error; err != nil {
		tx.Rollback()
		return err
	}

	analytics := &models.Analytics{
		SuccessfulRegisters:            registers,
		SuccessfulREgisterLastUpdateAt: time.Now(),
	}
	if err := tx.Create(analytics).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func (r *analyticsRepo) GetQuantityLogins() (int32, error) {
	var totalLogins int32
	result := r.db.Model(&models.Analytics{}).Select("successful_logins").Scan(&totalLogins)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalLogins, nil
}

func (r *analyticsRepo) GetQuantityRegister() (int32, error) {
	var totalRegistrations int32
	result := r.db.Model(&models.Analytics{}).Select("successful_register").Scan(&totalRegistrations)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalRegistrations, nil
}
