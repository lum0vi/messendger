package service

import (
	"gateway/internal/models"
	"gateway/internal/repository"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	return s.repo.Auth.Register(req)
}

func (s *AuthService) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	return s.repo.Auth.Login(req)
}
