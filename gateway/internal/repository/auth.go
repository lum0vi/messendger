package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	customErr "gateway/internal/errors"
	"gateway/internal/models"
	"io"
	"net/http"
)

type AuthRepository struct {
	cfg *config.Config
}

func NewAuthRepository(cfg *config.Config) *AuthRepository {
	return &AuthRepository{
		cfg: cfg,
	}
}

func (r *AuthRepository) Register(req *models.RegisterRequest) (*models.RegisterResponse, error) {
	jsReq, err := json.Marshal(req)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not marshal auth request: %v", err))
	}
	res, err := http.Post("http://"+r.cfg.AuthServiceHost+":"+r.cfg.AuthServicePort+"/auth/register", "application/json", bytes.NewBuffer(jsReq))
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not register auth request: %v", err))
	}
	defer res.Body.Close()

	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not read response body: %v", err))
	}

	if res.StatusCode == http.StatusConflict {
		return nil, customErr.NewCustomError(res.StatusCode, "already exists")
	}
	if res.StatusCode != http.StatusOK {
		return nil, customErr.NewCustomError(res.StatusCode, fmt.Sprintf("could not register auth request: %v", err))
	}

	var jsResp models.RegisterResponse
	err = json.Unmarshal(resp, &jsResp)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not marshal auth response: %v", err))
	}
	return &jsResp, nil
}

func (r *AuthRepository) Login(req *models.LoginRequest) (*models.LoginResponse, error) {
	jsReq, err := json.Marshal(req)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not marshal auth request: %v", err))
	}
	res, err := http.Post("http://"+r.cfg.AuthServiceHost+":"+r.cfg.AuthServicePort+"/auth/login", "application/json", bytes.NewBuffer(jsReq))
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not register auth request: %v", err))
	}
	defer res.Body.Close()
	resp, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not read response body: %v", err))
	}
	if res.StatusCode != http.StatusOK {
		return nil, customErr.NewCustomError(res.StatusCode, fmt.Sprintf("could not register auth request: %v", err))
	}

	var jsResp models.LoginResponse
	err = json.Unmarshal(resp, &jsResp)
	if err != nil {
		return nil, customErr.NewCustomError(500, fmt.Sprintf("could not marshal auth response: %v", err))
	}
	return &jsResp, nil
}
