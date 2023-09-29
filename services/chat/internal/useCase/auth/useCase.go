package auth

import (
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/security"
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/repository/storage/postgres"
	"time"

	"github.com/google/uuid"
	"github.com/thanhpk/randstr"
)

type authUC struct {
	cfg      *config.Config
	authRepo postgres.AuthRepository
	log      *logger.ZapLogger
}

func New(cfg *config.Config, authRepo postgres.AuthRepository, log *logger.ZapLogger) *authUC {
	return &authUC{cfg: cfg, authRepo: authRepo, log: log}
}

func (uc *authUC) convertToResponse(user *models.User) *models.UserResponse {
	return &models.UserResponse{
		ID:                user.ID,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Email:             user.Email,
		Phone:             user.Phone,
		ProfilePictureURL: user.ProfilePictureURL,
		City:              user.City,
		Gender:            user.Gender,
		Birthday:          user.Birthday,
		Age:               user.Age,
		CreatedAt:         user.CreatedAt,
		UpdatedAt:         user.UpdatedAt,
		Role:              user.Role,
		LastLoginAt:       user.LastLoginAt,
	}
}

func (uc *authUC) Register(user *models.User) (*models.UserResponse, error) {
	existsUserByEmail, err := uc.authRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if existsUserByEmail != nil {
		return nil, errors.New(e.ErrEmailNotExists)
	}

	existsUserByPhone, err := uc.authRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, err
	}

	if existsUserByPhone != nil {
		return nil, errors.New(e.ErrPhoneNotExists)
	}

	if err := models.BeforeCreate(user); err != nil {
		uc.log.Error("Error in BeforeCreate:", err)
		return nil, err
	}

	newUser, err := uc.authRepo.Register(user)
	if err != nil {
		uc.log.Error("Error in Register:", err)
		return nil, err
	}

	response := uc.convertToResponse(newUser)
	return response, nil
}

func (uc *authUC) Login(user *models.User) (*models.UserResponse, error) {
	var foundUser *models.User

	existsUserByEmail, err := uc.authRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, errors.New(e.ErrInvalidCredentials)
	}

	existsUserByPhone, err := uc.authRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, errors.New(e.ErrInvalidCredentials)
	}

	if existsUserByEmail != nil {
		foundUser = existsUserByEmail
	} else if existsUserByPhone != nil {
		foundUser = existsUserByPhone
	} else {
		return nil, errors.New(e.ErrInvalidCredentials)
	}

	if err := security.ComparePasswords(foundUser.Password, user.Password); err != nil {
		return nil, errors.New(e.ErrInvalidCredentials)
	}

	err = uc.authRepo.UpdateLastLogin(foundUser.ID)
	if err != nil {
		return nil, err
	}

	response := uc.convertToResponse(foundUser)
	return response, nil
}

func (uc *authUC) GetUsers() ([]*models.UserResponse, error) {
	users, err := uc.authRepo.GetUsers()
	if err != nil {
		uc.log.Error("Error in GetUsers:", err)
		return nil, err
	}

	responseUsers := make([]*models.UserResponse, 0, len(users))

	for _, user := range users {
		convertUser := uc.convertToResponse(user)
		responseUsers = append(responseUsers, convertUser)
	}

	return responseUsers, nil
}

func (uc *authUC) GetUserByID(userID uuid.UUID) (*models.UserResponse, error) {
	user, err := uc.authRepo.GetByID(userID)
	if err != nil {
		uc.log.Error("Error in GetUser:", err)
		return nil, err
	}

	response := uc.convertToResponse(user)
	return response, nil
}

func (uc *authUC) UpdateUser(user *models.User, userID uuid.UUID) (*models.UserResponse, error) {
	existingUser, err := uc.authRepo.GetByID(userID)
	if err != nil {
		return nil, errors.New(e.ErrUserNotFound)
	}

	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Email = user.Email
	existingUser.Role = user.Role
	existingUser.ProfilePictureURL = user.ProfilePictureURL
	existingUser.Phone = user.Phone
	existingUser.City = user.City
	existingUser.Gender = user.Gender
	existingUser.Birthday = user.Birthday
	existingUser.Age = user.Age

	updatedUser, err := uc.authRepo.UpdateUser(existingUser)
	if err != nil {
		uc.log.Error("Error in UpdateUser:", err)
		return nil, err
	}

	response := uc.convertToResponse(updatedUser)
	return response, nil
}

func (uc *authUC) DeleteUser(userID uuid.UUID) error {
	err := uc.authRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *authUC) ForgotPassword(email string) (*models.UserResponse, string, error) {
	existsUser, err := uc.authRepo.GetByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if existsUser == nil {
		return nil, "", errors.New(e.ErrEmailNotExists)
	}

	resetToken := randstr.String(20)
	passwordResetToken, err := utils.Encode(resetToken)
	if err != nil {
		return nil, "", err
	}
	existsUser.PasswordResetToken = passwordResetToken
	existsUser.PasswordResetAt = time.Now().Add(time.Minute * 15)

	updatedUser, err := uc.authRepo.UpdateUser(existsUser)
	if err != nil {
		return nil, "", err
	}

	response := uc.convertToResponse(updatedUser)
	return response, resetToken, nil
}

func (uc *authUC) ResetPassword(token, password string) error {
	existsUser, err := uc.authRepo.GetByToken(token)
	if err != nil {
		return err
	}

	existsUser.Password = password
	existsUser.PasswordResetToken = ""

	_, err = uc.authRepo.UpdateUser(existsUser)
	if err != nil {
		return err
	}

	return nil
}
