package repository

import (
	"auth/internal/errors"
	"auth/internal/hash"
	"auth/internal/models"
	"fmt"
	"net/http"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type PostgresAuthRepository struct {
	db *sqlx.DB
}

func NewPostgresAuthRepository(db *sqlx.DB) *PostgresAuthRepository {
	return &PostgresAuthRepository{
		db: db,
	}
}

func (r *PostgresAuthRepository) Register(req *models.RegisterRequest) (string, error) {
	fmt.Println(322332)
	query := "INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id"
	var id string
	err := r.db.QueryRow(query, req.Username, req.Password, req.Email).Scan(&id)
	if err != nil {
		// Преобразуем к pq.Error
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" && pgErr.Constraint == "users_username_key" {
				return "", errors.NewHttpError(http.StatusConflict, "username already in use")
			}
		}
		// Фолбэк
		return "", errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}

	return id, err
}

func (r *PostgresAuthRepository) Login(req *models.LoginRequest) (string, error) {
	query := "SELECT id, password_hash FROM users WHERE username = $1"
	var id, passHash string
	err := r.db.QueryRow(query, req.Username).Scan(&id, &passHash)
	if err != nil {
		return "", errors.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	ok := hash.VerifyPass(req.Password, passHash)
	if !ok {
		return "", errors.NewHttpError(http.StatusUnauthorized, "invalid password")
	}
	return id, err
}
