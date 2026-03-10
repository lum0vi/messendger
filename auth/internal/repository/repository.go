package repository

import (
	"auth/internal/models"

	"github.com/jmoiron/sqlx"
)

type PostgresAuth interface {
	Register(req *models.RegisterRequest) (string, error)
	Login(req *models.LoginRequest) (string, error)
}
type Repository struct {
	PostgresAuth PostgresAuth
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		PostgresAuth: NewPostgresAuthRepository(db),
	}
}
