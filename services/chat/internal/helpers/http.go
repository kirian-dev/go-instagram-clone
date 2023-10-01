package helpers

import "go-instagram-clone/services/chat/internal/domain/models"

func ConvertToResponseUser(user *models.User) *models.UserResponse {
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
		IsVerify:          user.IsVerify,
	}
}
