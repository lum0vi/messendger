package repository

import (
	"fmt"
	"net/http"
	"strings"
	"user/internal/errors"
	"user/internal/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) GetMe(id string) (*models.GetMeResponse, error) {
	query := "SELECT id, username, password_hash, email, created_at, updated_at FROM users WHERE id = $1"
	var user models.GetMeResponse
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Password,
		&user.Email,
		&user.CreateAt,
		&user.UpdateAt,
	)
	fmt.Println(user)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusNotFound, fmt.Sprintf("user not found: %s", err))
	}
	return &user, nil
}

func (r *UserRepository) UpdateMe(id string, req *models.UpdateMeRequest) error {
	setParts := []string{}
	args := []interface{}{}
	argIdx := 1

	if req.Username != "" {
		setParts = append(setParts, fmt.Sprintf("username = $%d", argIdx))
		args = append(args, req.Username)
		argIdx++
	}
	if req.Email != "" {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIdx))
		args = append(args, req.Email)
		argIdx++
	}
	if req.Password != "" {
		setParts = append(setParts, fmt.Sprintf("password_hash = $%d", argIdx))
		args = append(args, req.Password)
		argIdx++
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = $%d", strings.Join(setParts, ", "), argIdx)
	args = append(args, id)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("failed to update user: %s", err))
	}

	return nil
}

func (r *UserRepository) GetUsers() (*models.GetUsersResponse, error) {
	var users []*models.UserForGetUsers
	query := "SELECT id, username, email, created_at, updated_at FROM users"
	err := r.db.Select(&users, query)
	if err != nil {
		return nil, err
	}

	return &models.GetUsersResponse{
		users,
	}, nil
}

func (r *UserRepository) GetUserByID(req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error) {
	query := "SELECT * FROM users WHERE id = $1"
	var res models.GetUserByIDResponse
	err := r.db.QueryRow(query, req.ID).Scan(
		&res.ID,
		&res.Username,
		&res.Email,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusNotFound, fmt.Sprintf("user not found: %s", err))
	}
	return &res, nil
}

func (r *UserRepository) GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error) {
	query := "SELECT * FROM users WHERE username = $1"
	var res models.GetUserByUsernameResponse
	err := r.db.QueryRow(query, req.Username).Scan(
		&res.ID,
		&res.Username,
		&res.Email,
		&res.Password,
		&res.CreatedAt,
		&res.UpdatedAt,
	)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusNotFound, fmt.Sprintf("user not found: %s", err))
	}
	return &res, nil
}
