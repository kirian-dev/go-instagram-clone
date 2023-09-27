package postgres

import (
	"go-instagram-clone/services/analytics/internal/models"
	"time"

	"gorm.io/gorm"
)

type AnalyticsRepository interface {
	SaveSuccessfulLogin(email, phone string) error
	SaveSuccessfulRegister(email, phone string) error
	GetQuantityLogins() (int32, error)
	GetQuantityRegister() (int32, error)
}

type analyticsRepo struct {
	db *gorm.DB
}

func NewAnalyticsRepository(db *gorm.DB) *analyticsRepo {
	return &analyticsRepo{db: db}
}

func (r *analyticsRepo) SaveSuccessfulLogin(email, phone string) error {
	analytics := &models.Analytics{}
	result := r.db.Where("email = ? OR phone = ?", email, phone).First(analytics)

	if result.Error != nil {
		return result.Error
	}
	if err := r.db.Model(analytics).Update("successful_logins", gorm.Expr("successful_logins + ?", 1)).Error; err != nil {
		return err
	}

	if err := r.db.Model(analytics).Update("successful_login_last_update_at", time.Now()).Error; err != nil {
		return err
	}

	return nil
}

func (r *analyticsRepo) SaveSuccessfulRegister(email, phone string) error {
	analytics := &models.Analytics{
		Email:                          email,
		Phone:                          phone,
		SuccessfulRegister:             1,
		SuccessfulRegisterLastUpdateAt: time.Now(),
		SuccessfulLogins:               0,
		SuccessfulLoginLastUpdateAt:    time.Time{},
	}

	return r.db.Create(analytics).Error
}

func (r *analyticsRepo) GetQuantityLogins() (int32, error) {
	var totalLogins int32
	result := r.db.Model(&models.Analytics{}).Select("SUM(successful_logins)").Scan(&totalLogins)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalLogins, nil
}

func (r *analyticsRepo) GetQuantityRegister() (int32, error) {
	var totalRegistrations int32
	result := r.db.Model(&models.Analytics{}).Select("SUM(successful_register)").Scan(&totalRegistrations)
	if result.Error != nil {
		return 0, result.Error
	}
	return totalRegistrations, nil
}
