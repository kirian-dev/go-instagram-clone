package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FileImport struct {
	ID                 uuid.UUID        `json:"id" gorm:"primaryKey"`
	StartTime          time.Time        `json:"start_time"`
	EndTime            time.Time        `json:"end_time"`
	Status             FileImportStatus `json:"status"`
	SuccessfulAccounts int              `json:"successful_accounts"`
	PendingAccounts    int              `json:"pending_accounts"`
	FailedAccounts     int              `json:"failed_accounts"`
}

type FileImportStatus string

const (
	ImportStatusInProgress FileImportStatus = "in_progress"
	ImportStatusSuccess    FileImportStatus = "success"
	ImportStatusFailed     FileImportStatus = "failed"
)

func (fi *FileImport) BeforeCreate(tx *gorm.DB) (err error) {
	fi.ID = uuid.New()
	return nil
}
