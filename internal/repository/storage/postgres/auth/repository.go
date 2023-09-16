package auth

import (
	"context"
	"database/sql"
	"go-instagram-clone/internal/domain/models"
	"go-instagram-clone/pkg/security"

	"github.com/google/uuid"
)

type authRepository struct {
	db *sql.DB
}

func NewAuthRepository(db *sql.DB) *authRepository {
	return &authRepository{db: db}
}

func (r *authRepository) Register(ctx context.Context, user *models.User) (*models.User, error) {
	var u models.User
	if err := r.db.QueryRowContext(ctx, createUserQuery, &user.FirstName, &user.LastName, &user.Email,
		&user.Password, &user.Role, &user.ProfilePictureURL, &user.Phone, &user.City,
		&user.Gender, &user.Birthday, &user.Age).Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.Phone,
		&u.ProfilePictureURL,
		&u.City,
		&u.Birthday,
		&u.Age,
		&u.Gender,
		&u.CreatedAt,
		&u.UpdatedAt,
		&u.LastLoginAt,
		&u.Role); err != nil {
		return nil, err
	}

	return &u, nil
}

func (r *authRepository) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	return r.getUserByQuery(ctx, getByEmailQuery, email)
}

func (r *authRepository) GetByPhone(ctx context.Context, phone string) (*models.User, error) {
	return r.getUserByQuery(ctx, getByPhoneQuery, phone)
}

func (r *authRepository) GetByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	return r.getUserByQuery(ctx, getByIDQuery, userID)
}

func (r *authRepository) GetUsers(ctx context.Context) ([]*models.User, error) {
	rows, err := r.db.QueryContext(ctx, getUsersQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*models.User
	for rows.Next() {
		user, err := r.scanUser(rows)
		if err != nil {
			return nil, err
		}
		security.DeletePassword(&user.Password)
		users = append(users, user)
	}

	return users, nil
}

func (r *authRepository) UpdateUser(ctx context.Context, user *models.User) (*models.User, error) {
	_, err := r.db.ExecContext(
		ctx,
		updateUserQuery,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Password,
		user.Role,
		user.ProfilePictureURL,
		user.Phone,
		user.City,
		user.Gender,
		user.Birthday,
		user.Age,
		user.ID,
	)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *authRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, deleteUserQuery, userID)
	if err != nil {
		return err
	}

	return nil
}

func (r *authRepository) UpdateLastLogin(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, updateLastLoginQuery, userID)
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) scanUser(rows *sql.Rows) (*models.User, error) {
	var user models.User

	if err := rows.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.Phone,
		&user.ProfilePictureURL,
		&user.City,
		&user.Birthday,
		&user.Age,
		&user.Gender,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.LastLoginAt,
		&user.Role,
	); err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *authRepository) getUserByQuery(ctx context.Context, query string, args ...interface{}) (*models.User, error) {
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if !rows.Next() {
		return nil, nil
	}

	return r.scanUser(rows)
}
