package fileImport

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/services/chat/internal/useCase"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type importHandlers struct {
	cfg          *config.Config
	log          *logger.ZapLogger
	fileImportUC useCase.FileImportUseCase
}

const (
	maxFileSize = 2 * 1024 * 1024
)

func New(cfg *config.Config, log *logger.ZapLogger, fileImportUC useCase.FileImportUseCase) *importHandlers {
	return &importHandlers{cfg: cfg, log: log, fileImportUC: fileImportUC}
}

// @Summary Get import file status
// @Description Get status about uploaded users from file
// @Accept json
// @Produce json
// @Tags FileImport
// @Param fileId path int true "fileId"
// @Success 200 {object} models.FileImport
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /import/files/{fileId}/upload [get]
func (h *importHandlers) GetImportStatus() echo.HandlerFunc {
	return func(c echo.Context) error {
		fileIdStr := c.Param("fileId")

		fileID, err := uuid.Parse(fileIdStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		file, err := h.fileImportUC.GetImportFileByID(fileID)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: e.ErrInternalServer})
		}

		return c.JSON(http.StatusOK, file)
	}
}

// @Summary Get import files
// @Description Get all import files
// @Accept json
// @Produce json
// @Tags FileImport
// @Success 200 {object} models.FileImport
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /import/files [get]
func (h *importHandlers) GetImportList() echo.HandlerFunc {
	return func(c echo.Context) error {
		files, err := h.fileImportUC.GetImportFiles()
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: e.ErrInternalServer})
		}

		return c.JSON(http.StatusOK, files)
	}
}

// @Summary Upload file
// @Description Upload file from file csv for uploaded users list
// @Accept json
// @Produce json
// @Tags Users
// @Success 201 {object} models.FileImport
// @Failure 400 {object} e.ErrorResponse "Bad Request"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /import/files [post]
func (h *importHandlers) UploadFile() echo.HandlerFunc {
	return func(c echo.Context) error {
		file, err := c.FormFile("file")
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: e.ErrFailedToGetFile})
		}

		if file.Size > maxFileSize {
			return c.JSON(http.StatusRequestEntityTooLarge, e.ErrorResponse{Error: e.ErrBigFileSize})
		}

		ext := filepath.Ext(file.Filename)
		if ext != ".csv" {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: e.ErrFileMustBeCSV})
		}

		go func() {
			err = h.fileImportUC.UploadFile(file)
			if err != nil {
				return
			}
		}()

		return c.JSON(http.StatusCreated, map[string]string{
			"message": "File successfully uploaded",
		})
	}
}
