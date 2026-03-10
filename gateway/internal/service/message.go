package service

import (
	"gateway/internal/models"
	"gateway/internal/repository"
)

type MessageService struct {
	repo *repository.Repository
}

func NewMessageService(repo *repository.Repository) *MessageService {
	return &MessageService{repo: repo}
}

func (s *MessageService) GetMessagesByChatID(chatID string) (*models.GetUserMessageResponse, error) {
	return s.repo.Message.GetMessagesByChatID(chatID)
}
