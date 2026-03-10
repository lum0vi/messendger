package service

import (
	"chat/internal/models"
	"chat/internal/repository"
)

type ChatService struct {
	repo *repository.Repository
}

func NewChatService(repo *repository.Repository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreatePrivateChat(userID string, req *models.CreatePrivateChatRequest) (string, error) {
	return s.repo.Chat.CreatePrivateChat(userID, req)
}

func (s *ChatService) CreatePublicChat(userID string, req *models.CreatePublicChatRequest) (string, error) {
	return s.repo.Chat.CreatePublicChat(userID, req)
}

func (s *ChatService) GetChats(id string) ([]string, error) {
	return s.repo.Chat.GetChats(id)
}

func (s *ChatService) GetUsersChat(chatID string) ([]string, error) {
	return s.repo.Chat.GetUsersChat(chatID)
}
