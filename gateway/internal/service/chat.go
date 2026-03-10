package service

import (
	"gateway/internal/models"
	"gateway/internal/repository"
)

type ChatService struct {
	repo *repository.Repository
}

func NewChatService(repo *repository.Repository) *ChatService {
	return &ChatService{
		repo: repo,
	}
}

func (s *ChatService) CreatePrivateChat(userID string, in *models.CreatePrivateChatRequest) (*models.CreatePrivateChatResponse, error) {
	return s.repo.Chat.CreatePrivateChat(userID, in)
}
func (s *ChatService) CreatePublicChat(userID string, in *models.CreatePublicChatRequest) (*models.CreatePublicChatResponse, error) {
	return s.repo.Chat.CreatePublicChat(userID, in)
}
func (s *ChatService) GetMeChats(userID string) (*models.GetMeChatsResponse, error) {
	return s.repo.Chat.GetMeChats(userID)
}
func (s *ChatService) GetChatUsers(UserID string, chatID string) (*models.GetChatUsersResponse, error) {
	return s.repo.Chat.GetChatUsers(UserID, chatID)
}
