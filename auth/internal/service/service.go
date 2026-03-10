package service

import (
	"auth/internal/models"
	"auth/internal/repository"
)

type Auth interface {
	Register(req *models.RegisterRequest) (string, error)
	Login(req *models.LoginRequest) (string, error)
}
type Service struct {
	Auth Auth
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth: NewAuthService(repo),
	}
}
