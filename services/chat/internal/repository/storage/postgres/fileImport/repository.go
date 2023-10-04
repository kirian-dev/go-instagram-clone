package chat

import (
	"go-instagram-clone/services/chat/internal/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type fileImportRepository struct {
	db *gorm.DB
}

func NewFileImportRepository(db *gorm.DB) *fileImportRepository {
	return &fileImportRepository{db: db}
}

func (r *fileImportRepository) CreateFile(fileImport *models.FileImport) (*models.FileImport, error) {
	if err := r.db.Create(fileImport).Error; err != nil {
		return nil, err
	}
	return fileImport, nil
}

func (r *fileImportRepository) UpdateFile(fileImport *models.FileImport) (*models.FileImport, error) {
	if err := r.db.Save(fileImport).Error; err != nil {
		return nil, err
	}
	return fileImport, nil
}

func (r *fileImportRepository) GetImportFiles() ([]*models.FileImport, error) {
	var files []*models.FileImport
	if err := r.db.Find(&files).Error; err != nil {
		return nil, err
	}
	return files, nil
}

func (r *fileImportRepository) GetImportFileByID(fileId uuid.UUID) (*models.FileImport, error) {
	var file *models.FileImport
	if err := r.db.Where("id = ?", fileId).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return file, nil
}
