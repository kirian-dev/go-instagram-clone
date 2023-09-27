package models

import "time"

type Analytics struct {
	SuccessfulLogins               int32     `json:"successful_logins"`
	SuccessfulLoginLastUpdateAt    time.Time `json:"successful_logins_update_at"`
	SuccessfulRegisters            int32     `json:"successful_register"`
	SuccessfulREgisterLastUpdateAt time.Time `json:"successful_register_update_at"`
}
