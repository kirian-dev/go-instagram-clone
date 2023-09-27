package models

import "time"

type Analytics struct {
	ID                             uint      `gorm:"primaryKey" json:"id"`
	Email                          string    `json:"email"`
	Phone                          string    `json:"phone"`
	SuccessfulLogins               int32     `json:"successful_logins"`
	SuccessfulLoginLastUpdateAt    time.Time `json:"successful_login_last_update_at"`
	SuccessfulRegister             int32     `json:"successful_register"`
	SuccessfulRegisterLastUpdateAt time.Time `json:"successful_register_update_at"`
}
