package auth

import (
	"context"
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/internal/repository/storage/postgres"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/pkg/security"

	"github.com/google/uuid"
)

type authUC struct {
	cfg      *config.Config
	authRepo postgres.AuthRepository
	log      *logger.ZapLogger
}

func New(cfg *config.Config, authRepo postgres.AuthRepository, log *logger.ZapLogger) *authUC {
	return &authUC{cfg: cfg, authRepo: authRepo, log: log}
}
func (uc *authUC) Register(ctx context.Context, user *models.User) (*models.User, error) {
	existsUserByEmail, err := uc.authRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, err
	}

	if existsUserByEmail != nil {
		return nil, errors.New(e.ErrEmailNotExists)
	}

	existsUserByPhone, err := uc.authRepo.GetByPhone(ctx, user.Phone)
	if err != nil {
		return nil, err
	}

	if existsUserByPhone != nil {
		return nil, errors.New(e.ErrPhoneNotExists)
	}

	if err := user.BeforeCreate(); err != nil {
		uc.log.Error("Error in BeforeCreate:", err)
		return nil, err
	}

	newUser, err := uc.authRepo.Register(ctx, user)
	if err != nil {
		uc.log.Error("Error in Register:", err)
		return nil, err
	}

	security.DeletePassword(&newUser.Password)

	return newUser, nil
}

func (uc *authUC) Login(ctx context.Context, user *models.User) (*models.User, error) {
	var foundUser *models.User

	existsUserByEmail, err := uc.authRepo.GetByEmail(ctx, user.Email)
	if err != nil {
		return nil, errors.New(e.ErrInvalidCredentials)
	}

	existsUserByPhone, err := uc.authRepo.GetByPhone(ctx, user.Phone)
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

	err = uc.authRepo.UpdateLastLogin(ctx, foundUser.ID)
	if err != nil {
		return nil, err
	}

	security.DeletePassword(&foundUser.Password)

	return foundUser, nil
}

func (uc *authUC) GetUsers(ctx context.Context) ([]*models.User, error) {
	users, err := uc.authRepo.GetUsers(ctx)
	if err != nil {
		uc.log.Error("Error in GetUsers:", err)
		return nil, err
	}

	return users, nil
}

func (uc *authUC) GetUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	user, err := uc.authRepo.GetByID(ctx, userID)
	if err != nil {
		uc.log.Error("Error in GetUser:", err)
		return nil, err
	}

	security.DeletePassword(&user.Password)
	return user, nil
}

func (uc *authUC) UpdateUser(
	ctx context.Context,
	user *models.User,
	userID uuid.UUID,
) (*models.User, error) {
	existingUser, err := uc.authRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New(e.ErrUserNotFound)
	}

	existingUser.FirstName = user.FirstName
	existingUser.LastName = user.LastName
	existingUser.Email = user.Email
	existingUser.Password = user.Password
	existingUser.Role = user.Role
	existingUser.ProfilePictureURL = user.ProfilePictureURL
	existingUser.Phone = user.Phone
	existingUser.City = user.City
	existingUser.Gender = user.Gender
	existingUser.Birthday = user.Birthday
	existingUser.Age = user.Age

	updatedUser, err := uc.authRepo.UpdateUser(ctx, existingUser)
	if err != nil {
		uc.log.Error("Error in UpdateUser:", err)
		return nil, err
	}
	security.DeletePassword(&updatedUser.Password)

	return updatedUser, nil
}

func (uc *authUC) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := uc.authRepo.DeleteUser(ctx, userID)
	if err != nil {
		return err
	}

	return nil
}
