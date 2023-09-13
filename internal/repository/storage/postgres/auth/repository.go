package auth

import (
	"context"
	"database/sql"
	"go-instagram-clone/internal/domain/models"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	u := &models.User{}
	if err := r.db.QueryRowContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.ProfilePictureURL, &user.Phone, &user.City,
		&user.Gender, &user.Birthday, &user.Age,
	).Scan(u); err != nil {
		return nil, err
	}

	return u, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, user *models.User) (*models.User, error) {
	foundUser := &models.User{}
	if err := r.db.QueryRowContext(ctx, getByEmailQuery, user.Email).Scan(foundUser); err != nil {
		return nil, err
	}
	return foundUser, nil
}
