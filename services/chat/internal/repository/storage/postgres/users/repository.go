package users

import (
	"go-instagram-clone/pkg/utils"
	"go-instagram-clone/services/chat/internal/domain/models"
	"go-instagram-clone/services/chat/internal/helpers"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type usersRepository struct {
	db *gorm.DB
}

func NewUsersRepository(db *gorm.DB) *usersRepository {
	return &usersRepository{db: db}
}

func (r *usersRepository) convertUsers(users []*models.User) []*models.UserResponse {
	var usersConverted []*models.UserResponse

	for _, user := range users {
		userConverted := helpers.ConvertToResponseUser(user)
		usersConverted = append(usersConverted, userConverted)
	}

	return usersConverted
}

func (r *usersRepository) GetByID(userID uuid.UUID) (*models.User, error) {
	var user models.User
	if err := r.db.Where("id = ?", userID).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *usersRepository) GetUsers(pag *utils.PaginationQuery) (*models.UserListResponse, error) {
	offset := pag.GetOffset()
	limit := pag.GetLimit()

	var totalCount int64

	if err := r.db.Model(&models.User{}).
		Count(&totalCount).
		Error; err != nil {
		return nil, err
	}

	var users []*models.User

	query := r.db.Offset(offset).Limit(limit)
	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	usersResponse := r.convertUsers(users)

	response := &models.UserListResponse{
		Users:      usersResponse,
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pag.GetSize()),
		Page:       pag.GetPage(),
		Size:       pag.GetSize(),
		HasMore:    utils.GetHasMore(pag.GetPage(), totalCount, pag.GetSize()),
	}

	return response, nil
}

func (r *usersRepository) UpdateUser(user *models.User) (*models.User, error) {
	if err := r.db.Save(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (r *usersRepository) DeleteUser(userID uuid.UUID) error {
	if err := r.db.Delete(&models.User{}, userID).Error; err != nil {
		return err
	}
	return nil
}

func (r *usersRepository) SearchByQuery(query string, pag *utils.PaginationQuery) (*models.UserListResponse, error) {
	offset := pag.GetOffset()
	limit := pag.GetLimit()

	var totalCount int64

	if err := r.db.Model(&models.User{}).
		Where("first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%").
		Count(&totalCount).
		Error; err != nil {
		return nil, err
	}

	var users []*models.User

	if err := r.db.Model(&models.User{}).
		Where("first_name LIKE ? OR last_name LIKE ?", "%"+query+"%", "%"+query+"%").
		Offset(offset).
		Limit(limit).
		Find(&users).
		Error; err != nil {
		return nil, err
	}

	usersResponse := r.convertUsers(users)

	response := &models.UserListResponse{
		Users:      usersResponse,
		TotalCount: totalCount,
		TotalPages: utils.GetTotalPages(totalCount, pag.GetSize()),
		Page:       pag.GetPage(),
		Size:       pag.GetSize(),
		HasMore:    utils.GetHasMore(pag.GetPage(), totalCount, pag.GetSize()),
	}

	return response, nil
}
