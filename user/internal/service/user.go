package service

import (
	"user/internal/hash"
	"user/internal/models"
	"user/internal/repository"
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
	if req.Password != "" {
		hashPass, err := hash.HashPass(req.Password)
		if err != nil {
			return err
		}
		req.Password = hashPass
	}
	return s.repo.User.UpdateMe(id, req)
}

func (s *UserService) GetUsers() (*models.GetUsersResponse, error) {
	return s.repo.User.GetUsers()
}

func (s *UserService) GetUserByID(req *models.GetUserByIDRequest) (*models.GetUserByIDResponse, error) {
	return s.repo.User.GetUserByID(req)
}

func (s *UserService) GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error) {
	return s.repo.User.GetUserByUsername(req)
}
