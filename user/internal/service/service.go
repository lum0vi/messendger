package service

import (
	"user/internal/models"
	"user/internal/repository"
)

type User interface {
	GetMe(id string) (*models.GetMeResponse, error)
	UpdateMe(id string, req *models.UpdateMeRequest) error
	GetUsers() (*models.GetUsersResponse, error)
	GetUserByID(req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error)
	GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error)
}

type Service struct {
	User User
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		User: NewUserService(repo),
	}
}
