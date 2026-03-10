package service

import (
	"message/internal/models"
	"message/internal/repository"
)

type Message interface {
	GetUserMessages(userID string) ([]*models.Message, error)
	GetMessagesByChatID(chatID string) ([]*models.Message, error)
}

type Service struct {
	Message Message
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Message: NewMessageService(repo),
	}
}
