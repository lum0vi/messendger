package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"gateway/internal/config"
	"gateway/internal/errors"
	"gateway/internal/models"
)

type UserRepository struct {
	cfg *config.Config
}

func NewUserRepository(cfg *config.Config) *UserRepository {
	return &UserRepository{cfg: cfg}
}

func (r *UserRepository) GetMe(id string) (*models.GetMeResponse, error) {
	req, err := http.NewRequest("GET", "http://"+r.cfg.UserServiceHost+":"+r.cfg.UserServicePort+"/me", nil)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id", id)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %w", err))
	}

	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.NewCustomError(http.StatusNotFound, "user not found")
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(resp.StatusCode, string(body))
	}
	var user models.GetMeResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal response: %w", err))
	}
	return &user, nil
}

func (r *UserRepository) UpdateMe(id string, req *models.UpdateMeRequest) error {
	jsReq, err := json.Marshal(req)
	if err != nil {
		return errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not marshal request: %w", err))
	}
	request, err := http.NewRequest("PUT", "http://"+r.cfg.UserServiceHost+":"+r.cfg.UserServicePort+"/me", bytes.NewBuffer(jsReq))
	if err != nil {
		return errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %w", err))
	}
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("id", id)
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		return errors.NewCustomError(resp.StatusCode, string(body))
	}
	return nil
}

func (r *UserRepository) GetUsers(id string) (*models.GetUsersResponse, error) {
	req, err := http.NewRequest("GET", "http://"+r.cfg.UserServiceHost+":"+r.cfg.UserServicePort+"/users", nil)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %w", err))
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id", id)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %w", err))
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(resp.StatusCode, string(body))
	}
	var users models.GetUsersResponse
	err = json.Unmarshal(body, &users)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal response: %w", err))
	}
	return &users, nil
}

func (r *UserRepository) GetUserById(req *models.GetUserByIdRequest) (*models.GetUserByIdResponse, error) {
	jsReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not marshal request: %w", err))
	}
	request, err := http.NewRequest("POST", "http://"+r.cfg.UserServiceHost+":"+r.cfg.UserServicePort+"/user/id", bytes.NewBuffer(jsReq))
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %w", err))
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(resp.StatusCode, string(body))
	}
	var user models.GetUserByIdResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal response: %w", err))
	}
	return &user, nil
}

func (r *UserRepository) GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error) {
	jsReq, err := json.Marshal(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not marshal request: %w", err))
	}
	request, err := http.NewRequest("POST", "http://"+r.cfg.UserServiceHost+":"+r.cfg.UserServicePort+"/user/name", bytes.NewBuffer(jsReq))
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %w", err))
	}
	request.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %w", err))
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %w", err))
	}
	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(resp.StatusCode, string(body))
	}
	var user models.GetUserByUsernameResponse
	err = json.Unmarshal(body, &user)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal response: %w", err))
	}
	return &user, nil
}
