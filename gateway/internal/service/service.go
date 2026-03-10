package service

import (
	"gateway/internal/models"
	"gateway/internal/repository"
)

type Auth interface {
	Register(req *models.RegisterRequest) (*models.RegisterResponse, error)
	Login(req *models.LoginRequest) (*models.LoginResponse, error)
}

type User interface {
	GetMe(id string) (*models.GetMeResponse, error)
	UpdateMe(id string, req *models.UpdateMeRequest) error
	GetUsers(id string) (*models.GetUsersResponse, error)
	GetUserById(req *models.GetUserByIdRequest) (*models.GetUserByIdResponse, error)
	GetUserByUsername(req *models.GetUserByUsernameRequest) (*models.GetUserByUsernameResponse, error)
}

type Chat interface {
	CreatePrivateChat(userID string, in *models.CreatePrivateChatRequest) (*models.CreatePrivateChatResponse, error)
	CreatePublicChat(userID string, in *models.CreatePublicChatRequest) (*models.CreatePublicChatResponse, error)
	GetMeChats(userID string) (*models.GetMeChatsResponse, error)
	GetChatUsers(UserID string, chatID string) (*models.GetChatUsersResponse, error)
}
type Message interface {
	GetMessagesByChatID(chatID string) (*models.GetUserMessageResponse, error)
}

type Service struct {
	Auth    Auth
	User    User
	Chat    Chat
	Message Message
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Auth:    NewAuthService(repo),
		User:    NewUserService(repo),
		Chat:    NewChatService(repo),
		Message: NewMessageService(repo),
	}
}
