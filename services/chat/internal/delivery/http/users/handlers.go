package users

import (
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"io"
	"os"
	"path/filepath"

	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/useCase"
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
)

type usersHandlers struct {
	cfg     *config.Config
	log     *logger.ZapLogger
	usersUC useCase.UsersUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, usersUC useCase.UsersUseCase) *usersHandlers {
	return &usersHandlers{cfg: cfg, log: log, usersUC: usersUC}
}

// @Summary Get all users
// @Description Get users
// @Tags Users
// @Accept  json
// @Produce  json
// @Success 200 {object} models.UsersListResponse
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /users/all [get]
func (h *usersHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

		pag, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		usersPag, err := h.usersUC.GetUsers(pag)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: e.ErrInternalServer})
		}

		return c.JSON(http.StatusOK, usersPag)
	}
}

// @Summary Get user by id
// @Description Get user
// @Tags Users
// @Accept  json
// @Produce  json
// @Param userId path int true "userId"
// @Success 200 {object} models.User
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /users/{id} [get]
func (h *usersHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		userIDStr := c.Param("userId")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		user, err := h.usersUC.GetUserByID(userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: e.ErrUserNotFound})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// @Summary Update user
// @Description The user can update himself
// @Tags Users
// @Accept json
// @Param userId path int true "userId"
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} e.ErrorResponse "Invalid body parameters"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /users/{id} [put]
func (h *usersHandlers) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		userID, err := utils.AuthorizeUser(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}
		user := &models.User{}

		if err := utils.ReadRequest(c, user); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		updatedUser, err := h.usersUC.UpdateUser(user, userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// @Summary Delete User
// @Description The user can delete himself
// @Tags Users
// @Param userId path int true "userId"
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /users/{userId} [delete]
func (h *usersHandlers) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		userID, err := utils.AuthorizeUser(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		err = h.usersUC.DeleteUser(userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusNoContent, nil)
	}
}

// @Summary Get user account
// @Description Get the account details of the authenticated user.
// @Accept json
// @Produce json
// @Tags Users
// @Param userId path int true "userId"
// @Success 200 {object} models.User
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Router /users/account [get]
func (h *usersHandlers) GetAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims, ok := c.Get("userClaims").(*utils.CustomClaims)
		if !ok {
			h.log.Error(e.ErrUserContextNotFound)
			return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrUnauthorized})
		}

		user, err := h.usersUC.GetUserByID(userClaims.UserID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: e.ErrUserNotFound})
		}

		res := map[string]interface{}{
			"account": user,
		}
		return c.JSON(http.StatusOK, res)
	}
}

// @Summary Update avatar
// @Description Update avatar for user account
// @Accept json
// @Produce json
// @Tags Users
// @Param userId path int true "userId"
// @Success 200 {string} "Avatar updated successfully"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Router /users/{userId}/avatar [patch]
func (h *usersHandlers) UpdateAvatar() echo.HandlerFunc {
	return func(c echo.Context) error {
		userID, err := utils.AuthorizeUser(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}
		log.Info(userID)
		file, err := c.FormFile("avatar")
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		avatarFileName := uuid.New().String() + filepath.Ext(file.Filename)
		avatarPath := filepath.Join("public/uploads", avatarFileName)

		src, err := file.Open()
		if err != nil {
			log.Error(err)
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		defer src.Close()

		dst, err := os.Create(avatarPath)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: err.Error()})
		}

		err = h.usersUC.UpdateAvatar(userID, avatarPath)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"message": "Avatar updated successfully",
		})
	}
}

// @Summary Search user
// @Description search user by first name and last name using pagination
// @Accept json
// @Produce json
// @Tags Users
// @Success 200 {object} models.UsersListResponse
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Router /users/search [get]
func (h *usersHandlers) SearchUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		pag, err := utils.GetPaginationFromCtx(c)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err.Error())
		}

		users, err := h.usersUC.SearchByQuery(c.QueryParam("query"), pag)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, users)
	}
}
