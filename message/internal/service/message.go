package service

import (
	"message/internal/models"

	"message/internal/repository"
)

type MessageService struct {
	repo *repository.Repository
}

func NewMessageService(repo *repository.Repository) *MessageService {
	return &MessageService{repo: repo}
}

func (r *MessageService) GetMessagesByChatID(chatID string) ([]*models.Message, error) {
	return r.repo.Message.GetMessagesByChatID(chatID)
}
func (r *MessageService) GetUserMessages(userID string) ([]*models.Message, error) {
	return r.repo.Message.GetUserMessages(userID)
}
