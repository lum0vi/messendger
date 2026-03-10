package service

import (
	customErr "auth/internal/errors"
	"auth/internal/hash"
	"auth/internal/jwt"
	"auth/internal/models"
	"auth/internal/repository"
	"net/http"
)

type AuthService struct {
	repo *repository.Repository
}

func NewAuthService(repo *repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Register(req *models.RegisterRequest) (string, error) {
	hashPass, err := hash.HashPass(req.Password)
	if err != nil {
		return "", err
	}
	req.Password = hashPass
	id, err := s.repo.PostgresAuth.Register(req)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (s *AuthService) Login(req *models.LoginRequest) (string, error) {
	id, err := s.repo.PostgresAuth.Login(req)
	if err != nil {
		return "", err
	}
	token, err := jwt.CreateJWT(id)
	if err != nil {
		return "", customErr.NewHttpError(http.StatusInternalServerError, err.Error())
	}
	return token, nil
}
