package service

import (
	"gateway/internal/models"
	"gateway/internal/repository"
)

type UserService struct {
	repo *repository.Repository
}

func NewUserService(repo *repository.Repository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetMe(id string) (*models.GetMeResponse, error) {
	return s.repo.User.GetMe(id)
}

func (s *UserService) UpdateMe(id string, req *models.UpdateMeRequest) error {
	return s.repo.User.UpdateMe(id, req)
}

func (s *UserService) GetUsers(id string) (*models.GetUsersResponse, error) {
	return s.repo.User.GetUsers(id)
}

func (s *UserService) GetUserById(req *models.GetUserByIdRequest) (*models.GetUserByIdResponse, error) {
	return s.repo.User.GetUserById(req)
}

func (s *UserService) GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error) {
	return s.repo.User.GetUserByUsername(req)
}
