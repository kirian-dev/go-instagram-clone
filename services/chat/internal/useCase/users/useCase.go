package users

import (
	"errors"
	"go-instagram-clone/config"
	"go-instagram-clone/pkg/e"
	"go-instagram-clone/pkg/logger"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"
	"go-instagram-clone/services/chat/internal/repository/storage/postgres"

	"github.com/google/uuid"
)

type usersUC struct {
	cfg       *config.Config
	usersRepo postgres.UsersRepository
	log       *logger.ZapLogger
}

func New(cfg *config.Config, usersRepo postgres.UsersRepository, log *logger.ZapLogger) *usersUC {
	return &usersUC{cfg: cfg, usersRepo: usersRepo, log: log}
}

func (uc *usersUC) GetUsers() ([]*models.UserResponse, error) {
	users, err := uc.usersRepo.GetUsers()
	if err != nil {
		uc.log.Error("Error in GetUsers:", err)
		return nil, err
	}

	responseUsers := make([]*models.UserResponse, 0, len(users))

	for _, user := range users {
		convertUser := helpers.ConvertToResponseUser(user)
		responseUsers = append(responseUsers, convertUser)
	}

	return responseUsers, nil
}

func (uc *usersUC) GetUserByID(userID uuid.UUID) (*models.UserResponse, error) {
	user, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		uc.log.Error("Error in GetUser:", err)
		return nil, err
	}

	response := helpers.ConvertToResponseUser(user)
	return response, nil
}

func (uc *usersUC) UpdateUser(user *models.User, userID uuid.UUID) (*models.UserResponse, error) {
	existingUser, err := uc.usersRepo.GetByID(userID)
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

	updatedUser, err := uc.usersRepo.UpdateUser(existingUser)
	if err != nil {
		uc.log.Error("Error in UpdateUser:", err)
		return nil, err
	}

	response := helpers.ConvertToResponseUser(updatedUser)
	return response, nil
}

func (uc *usersUC) DeleteUser(userID uuid.UUID) error {
	err := uc.usersRepo.DeleteUser(userID)
	if err != nil {
		return err
	}

	return nil
}

func (uc *usersUC) UpdateAvatar(userID uuid.UUID, avatarPath string) error {
	existingUser, err := uc.usersRepo.GetByID(userID)
	if err != nil {
		return err
	}

	existingUser.ProfilePictureURL = avatarPath

	_, err = uc.usersRepo.UpdateUser(existingUser)
	if err != nil {
		return err
	}

	return nil
}
