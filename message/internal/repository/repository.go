package repository

import (
	"message/internal/models"

	"github.com/jmoiron/sqlx"
)

type Message interface {
	Save(msg *models.Message) (string, error)
	GetUserMessages(userID string) ([]*models.Message, error)
	GetMessagesByChatID(chatID string) ([]*models.Message, error)
	UsersSendMess(chatID string, senderID string) (*[]string, error)
}
type Repository struct {
	Message Message
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Message: NewMessageRepo(db),
	}
}
