package repository

import (
	"encoding/json"
	"fmt"
	"gateway/internal/config"
	"gateway/internal/errors"
	"gateway/internal/models"
	"io"
	"net/http"
)

type MessageRepository struct {
	cfg *config.Config
}

func NewMessageRepository(cfg *config.Config) *MessageRepository {
	return &MessageRepository{cfg: cfg}
}

func (r *MessageRepository) GetMessagesByChatID(chatID string) (*models.GetUserMessageResponse, error) {
	url := fmt.Sprintf("http://%s:%s/message/chat/%s", r.cfg.MessageServiceHost, r.cfg.MessageServicePort, chatID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not create request: %v", err))
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not send request: %v", err))
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not read response: %v", err))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.NewCustomError(resp.StatusCode, string(body))
	}

	var out models.GetUserMessageResponse
	if err := json.Unmarshal(body, &out); err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("could not unmarshal response: %v", err))
	}

	return &out, nil
}
