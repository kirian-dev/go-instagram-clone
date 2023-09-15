package auth

import (
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/useCase"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	cfg    *config.Config
	log    *logger.ZapLogger
	authUC useCase.AuthUseCase
}

func New(cfg *config.Config, log *logger.ZapLogger, authUC useCase.AuthUseCase) *authHandlers {

	return &authHandlers{cfg: cfg, log: log, authUC: authUC}
}

func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		if c.Request().Header.Get("Content-Type") != "application/json" {
			h.log.Error(e.ErrInvalidFormat)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.ErrInvalidFormat})
		}

		user := &models.User{}

		if err := utils.ReadRequest(c, user); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		createdUser, err := h.authUC.Register(ctx, user)

		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		accessToken, refreshToken, err := utils.GenerateJWTTokens(createdUser.Email, createdUser.Role, h.cfg)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrInternalServer)
		}

		res := map[string]interface{}{
			"account":      createdUser,
			"access_token": accessToken,
		}
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		return c.JSON(http.StatusCreated, res)
	}
}

func (h *authHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		type Login struct {
			Email    string `json:"email,omitempty" db:"email" validate:"omitempty,email,lte=60"`
			Password string `json:"password,omitempty" db:"password" validate:"required,omitempty,gte=6"`
			Phone    string `json:"phone,omitempty" db:"phone" validate:"omitempty,e164"`
		}
		login := &Login{}
		ctx := c.Request().Context()

		if c.Request().Header.Get("Content-Type") != "application/json" {
			h.log.Error(e.ErrInvalidFormat)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.ErrInvalidFormat})
		}

		if err := utils.ReadRequest(c, login); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		existsUser, err := h.authUC.Login(ctx, &models.User{
			Email:    login.Email,
			Password: login.Password,
			Phone:    login.Phone,
		})
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		accessToken, refreshToken, err := utils.GenerateJWTTokens(existsUser.Email, existsUser.Role, h.cfg)

		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrInternalServer)
		}

		res := map[string]interface{}{
			"account":      existsUser,
			"access_token": accessToken,
		}
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})
		return c.JSON(http.StatusOK, res)
	}
}

func (h *authHandlers) RefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		cookie, err := c.Cookie("refresh_token")
		if err != nil {
			return c.JSON(http.StatusUnauthorized, e.ErrRefreshNotFound)
		}
		refreshToken := cookie.Value

		claims, err := utils.ValidateJWTToken(refreshToken, h.cfg)
		if err != nil {
			return c.JSON(http.StatusUnauthorized, e.ErrInvalidRefreshToken)
		}

		expiresAtUnix := claims.ExpiresAt.Time.Unix()

		if expiresAtUnix < time.Now().Unix() {
			return c.JSON(http.StatusUnauthorized, e.ErrRefreshTokenExpired)
		}

		accessToken, refreshToken, err := utils.GenerateJWTTokens(claims.Email, claims.Role, h.cfg)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrInternalServer)
		}
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		res := map[string]interface{}{
			"access_token": accessToken,
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *authHandlers) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    "",
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
			MaxAge:   -1,
		})

		res := map[string]string{
			"message": "Logout successfully",
		}

		return c.JSON(http.StatusOK, res)
	}
}

func (h *authHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		users, err := h.authUC.GetUsers(ctx)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusInternalServerError, e.ErrorResponse{Error: e.ErrInternalServer})
		}

		if len(users) == 0 {
			return c.JSON(http.StatusOK, []interface{}{})
		}
		return c.JSON(http.StatusOK, users)
	}
}

func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userIDStr := c.Param("userId")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		user, err := h.authUC.GetUserByID(ctx, userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: e.ErrUserNotFound})
		}

		return c.JSON(http.StatusOK, user)
	}
}

func (h *authHandlers) UpdateUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userIDStr := c.Param("userId")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}

		user := &models.User{}

		if err := utils.ReadRequest(c, user); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		updatedUser, err := h.authUC.UpdateUser(ctx, user, userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}
func (h *authHandlers) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		userIDStr := c.Param("userId")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		err = h.authUC.DeleteUser(ctx, userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: err.Error()})
		}
		return nil
	}
}
