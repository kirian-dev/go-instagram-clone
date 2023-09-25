package auth

import (
	"context"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/utils"
	pb "go-instagram-clone/services/analytics/cmd/proto"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/useCase"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

type authHandlers struct {
	cfg             *config.Config
	log             *logger.ZapLogger
	authUC          useCase.AuthUseCase
	analyticsClient pb.AnalyticsServiceClient
}

func New(cfg *config.Config, log *logger.ZapLogger, authUC useCase.AuthUseCase, analyticsClient pb.AnalyticsServiceClient) *authHandlers {

	return &authHandlers{cfg: cfg, log: log, authUC: authUC, analyticsClient: analyticsClient}
}

// @Summary Register
// @Description Register a new user with the provided credentials, return access token and refresh token.
// @Accept json
// @Produce json
// @Tags Auth
// @Success 201 {object} map[string]interface{} "Response with access token and refresh token"
// @Failure 400 {object}  e.ErrorResponse "Invalid body parameters"
// @Router /auth/register [post]
func (h *authHandlers) Register() echo.HandlerFunc {
	return func(c echo.Context) error {

		if c.Request().Header.Get("Content-Type") != "application/json" {
			h.log.Error(e.ErrInvalidFormat)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.ErrInvalidFormat})
		}

		user := &models.User{}

		if err := utils.ReadRequest(c, user); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		createdUser, err := h.authUC.Register(user)

		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		accessToken, refreshToken, err := utils.GenerateJWTTokens(createdUser.Email, createdUser.Role, createdUser.ID, h.cfg)
		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrInternalServer)
		}

		res := map[string]interface{}{
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

// @Summary Login
// @Description Authenticate a user with their email or phone and password.
// @Accept json
// @Produce json
// @Tags Auth
// @Success 201 {object} map[string]interface{} "Response with access token and refresh token"
// @Failure 400 {object} e.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *authHandlers) Login() echo.HandlerFunc {
	return func(c echo.Context) error {

		type Login struct {
			Email    string `json:"email,omitempty" validate:"omitempty,email,lte=60"`
			Password string `json:"password,omitempty" validate:"required,omitempty,gte=6"`
			Phone    string `json:"phone,omitempty" validate:"omitempty,e164"`
		}
		login := &Login{}

		if c.Request().Header.Get("Content-Type") != "application/json" {
			h.log.Error(e.ErrInvalidFormat)
			return c.JSON(http.StatusBadRequest, map[string]string{"error": e.ErrInvalidFormat})
		}

		if err := utils.ReadRequest(c, login); err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		existsUser, err := h.authUC.Login(&models.User{
			Email:    login.Email,
			Password: login.Password,
			Phone:    login.Phone,
		})
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusBadRequest, e.ErrorResponse{Error: err.Error()})
		}

		accessToken, refreshToken, err := utils.GenerateJWTTokens(existsUser.Email, existsUser.Role, existsUser.ID, h.cfg)

		if err != nil {
			return c.JSON(http.StatusBadRequest, e.ErrInternalServer)
		}

		res := map[string]interface{}{
			"access_token": accessToken,
		}
		c.SetCookie(&http.Cookie{
			Name:     "refresh_token",
			Value:    refreshToken,
			HttpOnly: true,
			Secure:   true,
			Path:     "/",
		})

		req := &pb.LoginRequest{
			SuccessfulLogins: 1,
		}

		_, err = h.analyticsClient.RecordLogin(context.Background(), req)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to record login in analytics service"})
		}
		return c.JSON(http.StatusOK, res)
	}
}

// @Summary Refresh access token and refresh token
// @Description Refresh the access token using the refresh token.
// @Accept json
// @Produce json
// @Tags Auth
// @Success 201 {object} map[string]interface{} "Response with access token and refresh token"
// @Failure 400 {object} e.ErrorResponse "Refresh token is invalid"
// @Router /auth/refresh-token [post]
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

		accessToken, refreshToken, err := utils.GenerateJWTTokens(claims.Email, claims.Role, claims.UserID, h.cfg)
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

// @Summary Logout
// @Description Logout the user by clearing their refresh token.
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {string} string "Logout successful"
// @Router /auth/logout [post]
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

// @Summary Get all users
// @Description Get users
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} models.User
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /auth/all [get]
func (h *authHandlers) GetUsers() echo.HandlerFunc {
	return func(c echo.Context) error {

		users, err := h.authUC.GetUsers()
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

// @Summary Get user by id
// @Description Get user
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param userId path int true "userId"
// @Success 200 {object} models.User
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /auth/{id} [get]
func (h *authHandlers) GetUserByID() echo.HandlerFunc {
	return func(c echo.Context) error {

		userIDStr := c.Param("userId")

		userID, err := uuid.Parse(userIDStr)
		if err != nil {
			return c.JSON(http.StatusBadRequest, err)
		}
		user, err := h.authUC.GetUserByID(userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: e.ErrUserNotFound})
		}

		return c.JSON(http.StatusOK, user)
	}
}

// @Summary Update user
// @Description The user can update himself
// @Tags Auth
// @Accept json
// @Param userId path int true "userId"
// @Produce json
// @Success 200 {object} models.User
// @Failure 400 {object} e.ErrorResponse "Invalid body parameters"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Failure 500 {object} e.ErrorResponse "Internal Server Error"
// @Router /auth/{id} [put]
func (h *authHandlers) UpdateUser() echo.HandlerFunc {
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

		updatedUser, err := h.authUC.UpdateUser(user, userID)
		if err != nil {
			h.log.Error(err)
			return c.JSON(http.StatusNotFound, e.ErrorResponse{Error: err.Error()})
		}

		return c.JSON(http.StatusOK, updatedUser)
	}
}

// @Summary Delete User
// @Description The user can delete himself
// @Tags Auth
// @Param userId path int true "userId"
// @Produce json
// @Success 204 "No Content"
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Failure 403 {object} e.ErrorResponse "Forbidden"
// @Router /auth/{userId} [delete]
func (h *authHandlers) DeleteUser() echo.HandlerFunc {
	return func(c echo.Context) error {

		userID, err := utils.AuthorizeUser(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, err.Error())
		}

		err = h.authUC.DeleteUser(userID)
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
// @Tags Auth
// @Success 200 {object} models.User
// @Failure 401 {object} e.ErrorResponse "Unauthorized"
// @Router /auth/account [get]
func (h *authHandlers) GetAccount() echo.HandlerFunc {
	return func(c echo.Context) error {
		userClaims, ok := c.Get("userClaims").(*utils.CustomClaims)
		if !ok {
			h.log.Error(e.ErrUserContextNotFound)
			return c.JSON(http.StatusUnauthorized, e.ErrorResponse{Error: e.ErrUnauthorized})
		}

		user, err := h.authUC.GetUserByID(userClaims.UserID)
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
