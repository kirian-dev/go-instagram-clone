package fileImport

import (
	"go-instagram-clone/services/chat/internal/delivery/http"
	"go-instagram-clone/services/chat/internal/middleware"

	"github.com/labstack/echo/v4"
)

func MapImportRoutes(importGroup *echo.Group, h http.FileImportHandlers, mw *middleware.MiddlewareManager) {
	importGroup.Use(mw.AuthJWTMiddleware())
	importGroup.Use(mw.AdminAuthMiddleware())
	importGroup.GET("/files/:fileId", h.GetImportStatus())
	importGroup.GET("/files", h.GetImportList())
	importGroup.POST("/files", h.UploadFile())
}
