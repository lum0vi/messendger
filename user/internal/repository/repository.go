package repository

import (
	"user/internal/models"

	"github.com/jmoiron/sqlx"
)

type User interface {
	GetMe(id string) (*models.GetMeResponse, error)
	UpdateMe(id string, req *models.UpdateMeRequest) error
	GetUsers() (*models.GetUsersResponse, error)
	GetUserByID(req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error)
	GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error)
}

type Repository struct {
	User User
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User: NewUserRepository(db),
	}
}
