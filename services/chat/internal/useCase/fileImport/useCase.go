package fileImport

import (
	"bufio"
	"encoding/csv"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/logger"

	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"
	"go-instagram-clone/services/chat/internal/repository/storage/postgres"
	"io"
	"mime/multipart"
	"time"

	"github.com/google/uuid"
)

type fileImportUC struct {
	cfg            *config.Config
	log            *logger.ZapLogger
	fileImportRepo postgres.FileImportRepository
	usersRepo      postgres.UsersRepository
	authRepo       postgres.AuthRepository
}

func New(cfg *config.Config, log *logger.ZapLogger, fileImportRepo postgres.FileImportRepository, usersRepo postgres.UsersRepository, authRepo postgres.AuthRepository) *fileImportUC {
	return &fileImportUC{cfg, log, fileImportRepo, usersRepo, authRepo}
}

func (uc *fileImportUC) GetImportFiles() ([]*models.FileImport, error) {
	return uc.fileImportRepo.GetImportFiles()
}

func (uc *fileImportUC) GetImportFileByID(fileID uuid.UUID) (*models.FileImport, error) {
	return uc.fileImportRepo.GetImportFileByID(fileID)
}

func (uc *fileImportUC) UploadFile(file *multipart.FileHeader) error {
	rowCount, err := helpers.CountCSVRows(file)
	if err != nil {
		uc.log.Error("Failed to count CSV rows", err)
		fileImport := &models.FileImport{
			ID:                 uuid.New(),
			StartTime:          time.Now(),
			Status:             models.ImportStatusFailed,
			SuccessfulAccounts: 0,
			PendingAccounts:    rowCount,
			FailedAccounts:     0,
		}
		uc.fileImportRepo.CreateFile(fileImport)
		return err
	}

	fileImport := &models.FileImport{
		ID:                 uuid.New(),
		StartTime:          time.Now(),
		Status:             models.ImportStatusInProgress,
		SuccessfulAccounts: 0,
		PendingAccounts:    rowCount,
		FailedAccounts:     0,
	}

	if _, err := uc.fileImportRepo.CreateFile(fileImport); err != nil {
		uc.log.Error("Failed to create file import record", err)
		fileImport.Status = models.ImportStatusFailed
		uc.fileImportRepo.UpdateFile(fileImport)
		return err
	}

	fileSrc, err := file.Open()
	if err != nil {
		uc.log.Error("Failed to open file", err)
		fileImport.Status = models.ImportStatusFailed
		uc.fileImportRepo.UpdateFile(fileImport)
		return err
	}
	defer fileSrc.Close()

	reader := csv.NewReader(bufio.NewReader(fileSrc))
	pendingAccountsCh := make(chan int)

	go func() {
		for pendingCount := range pendingAccountsCh {
			if pendingCount < 0 {
				pendingCount = 0
			}
			fileImport.PendingAccounts = pendingCount
			if _, err := uc.fileImportRepo.UpdateFile(fileImport); err != nil {
				uc.log.Error("Failed to update file import status", err)
			}
		}
	}()

	for {
		line, err := reader.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			uc.log.Error("Error reading file line", err)
			fileImport.Status = models.ImportStatusFailed
			uc.fileImportRepo.UpdateFile(fileImport)
			return err
		}

		if line[2] == "email" {
			continue
		}

		email := line[2]

		existingUser, err := uc.authRepo.GetByEmail(email)
		if err != nil {
			uc.log.Error("Error checking for existing user", err)
			fileImport.Status = models.ImportStatusFailed
			uc.fileImportRepo.UpdateFile(fileImport)
			return err
		}

		if existingUser != nil {
			uc.log.Warn("User with the same email already exists", email)
			fileImport.FailedAccounts++
		} else {
			user := &models.User{
				ID:        uuid.New(),
				FirstName: line[0],
				LastName:  line[1],
				Email:     email,
			}

			_, err = uc.usersRepo.CreateUser(user)
			if err != nil {
				uc.log.Error("Error creating user", err)
				fileImport.FailedAccounts++
			} else {
				fileImport.SuccessfulAccounts++
			}
		}

		fileImport.PendingAccounts--
		pendingAccountsCh <- fileImport.PendingAccounts
	}

	close(pendingAccountsCh)

	fileImport.Status = models.ImportStatusSuccess
	fileImport.EndTime = time.Now()

	if _, err := uc.fileImportRepo.UpdateFile(fileImport); err != nil {
		uc.log.Error("Failed to update file import status", err)
		fileImport.Status = models.ImportStatusFailed
		uc.fileImportRepo.UpdateFile(fileImport)
		return err
	}

	return nil
}
