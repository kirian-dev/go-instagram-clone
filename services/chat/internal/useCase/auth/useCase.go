package auth

import (
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/security"
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"
	"go-instagram-clone/services/chat/internal/repository/storage/postgres"
	"time"

	"github.com/thanhpk/randstr"
)

type authUC struct {
	cfg       *config.Config
	authRepo  postgres.AuthRepository
	usersRepo postgres.UsersRepository
	log       *logger.ZapLogger
}

func New(cfg *config.Config, authRepo postgres.AuthRepository, usersRepo postgres.UsersRepository, log *logger.ZapLogger) *authUC {
	return &authUC{cfg: cfg, authRepo: authRepo, usersRepo: usersRepo, log: log}
}

func (uc *authUC) Register(user *models.User) (*models.UserResponse, error) {
	existsUserByEmail, err := uc.authRepo.GetByEmail(user.Email)
	if err != nil {
		return nil, err
	}

	if existsUserByEmail != nil {
		return nil, errors.New(e.ErrEmailExists)
	}

	existsUserByPhone, err := uc.authRepo.GetByPhone(user.Phone)
	if err != nil {
		return nil, err
	}

	if existsUserByPhone != nil {
		return nil, errors.New(e.ErrPhoneNotExists)
	}

	newUser, err := uc.authRepo.Register(user)
	if err != nil {
		uc.log.Error("Error in Register:", err)
		return nil, err
	}

	response := helpers.ConvertToResponseUser(newUser)
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

	response := helpers.ConvertToResponseUser(foundUser)
	return response, nil
}

func (uc *authUC) ForgotPassword(email string) (*models.UserResponse, string, error) {
	existsUser, err := uc.authRepo.GetByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if existsUser == nil {
		return nil, "", errors.New(e.ErrEmailNotExists)
	}

	if !existsUser.IsVerify {
		return nil, "", errors.New(e.ErrEmailMustBeVerified)
	}

	resetToken := randstr.String(20)
	passwordResetToken, err := utils.Encode(resetToken)
	if err != nil {
		return nil, "", err
	}
	existsUser.PasswordResetToken = passwordResetToken
	existsUser.PasswordResetAt = time.Now().Add(time.Minute * 15)

	updatedUser, err := uc.usersRepo.UpdateUser(existsUser)
	if err != nil {
		return nil, "", err
	}

	response := helpers.ConvertToResponseUser(updatedUser)
	return response, resetToken, nil
}

func (uc *authUC) ResetPassword(token, password string) error {
	existsUser, err := uc.authRepo.GetByToken(token)
	if err != nil {
		return err
	}

	existsUser.Password = password
	existsUser.PasswordResetToken = ""

	_, err = uc.usersRepo.UpdateUser(existsUser)
	if err != nil {
		return err
	}

	return nil
}

func (uc *authUC) SendVerificationEmail(email string) (*models.UserResponse, string, error) {
	existsUser, err := uc.authRepo.GetByEmail(email)
	if err != nil {
		return nil, "", err
	}

	if existsUser == nil {
		return nil, "", errors.New(e.ErrEmailNotExists)
	}

	verificationCode := randstr.String(20)
	existsUser.VerificationCode = verificationCode

	updatedUser, err := uc.usersRepo.UpdateUser(existsUser)
	if err != nil {
		return nil, "", err
	}

	return helpers.ConvertToResponseUser(updatedUser), verificationCode, nil
}

func (uc *authUC) VerifyEmail(code string) error {
	existsUser, err := uc.authRepo.GetByCode(code)
	if err != nil {
		return err
	}

	existsUser.VerificationCode = ""
	existsUser.IsVerify = true

	updated, err := uc.usersRepo.UpdateUser(existsUser)
	if err != nil {
		return err
	}
	uc.log.Info(updated)
	return nil
}
