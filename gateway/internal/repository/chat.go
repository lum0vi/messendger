package repository

import (
	"bytes"
	"encoding/json"
	"gateway/internal/config"
	"gateway/internal/errors"
	"gateway/internal/models"
	"io"
	"net/http"
)

type ChatRepository struct {
	cfg *config.Config
}

func NewChatRepository(cfg *config.Config) *ChatRepository {
	return &ChatRepository{cfg: cfg}
}

func (r *ChatRepository) CreatePrivateChat(userID string, in *models.CreatePrivateChatRequest) (*models.CreatePrivateChatResponse, error) {
	url := "http://" + r.cfg.ChatServiceHost + ":" + r.cfg.ChatServicePort + "/chat/private"
	jsonReq, err := json.Marshal(in)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	var response models.CreatePrivateChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	return &response, nil
}

func (r *ChatRepository) CreatePublicChat(userID string, in *models.CreatePublicChatRequest) (*models.CreatePublicChatResponse, error) {
	url := "http://" + r.cfg.ChatServiceHost + ":" + r.cfg.ChatServicePort + "/chat/public"
	jsonReq, err := json.Marshal(in)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonReq))
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("id", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	var response models.CreatePublicChatResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	return &response, nil
}

func (r *ChatRepository) GetMeChats(userID string) (*models.GetMeChatsResponse, error) {
	url := "http://" + r.cfg.ChatServiceHost + ":" + r.cfg.ChatServicePort + "/chat"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("id", userID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	var response models.GetMeChatsResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	return &response, nil
}

func (r *ChatRepository) GetChatUsers(UserID string, chatID string) (*models.GetChatUsersResponse, error) {
	url := "http://" + r.cfg.ChatServiceHost + ":" + r.cfg.ChatServicePort + "/chat/" + chatID + "/users"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	req.Header.Set("id", UserID)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	var response models.GetChatUsersResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, err.Error())
	}
	return &response, nil

}
