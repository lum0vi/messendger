package repository

import (
	"fmt"
	"message/internal/errors"
	"message/internal/models"
	"net/http"

	"github.com/jmoiron/sqlx"
)

type MessageRepo struct {
	db *sqlx.DB
}

func NewMessageRepo(db *sqlx.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

func (r *MessageRepo) Save(msg *models.Message) (string, error) {
	query := "INSERT INTO messages (chat_id, sender_id, content) VALUES ($1, $2, $3) RETURNING id, EXTRACT(EPOCH FROM sent_at)::bigint * 1000"
	var id string
	var sentAt int64
	err := r.db.QueryRow(query, msg.ChatID, msg.SenderID, msg.Content).Scan(&id, &sentAt)
	if err != nil {
		return id, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("failed to save message: %v", err))
	}
	msg.ID = id
	msg.SentAt = sentAt
	return id, nil
}

func (r *MessageRepo) GetUserMessages(userID string) ([]*models.Message, error) {
	var messages []models.Message

	query := `
		SELECT m.id, m.chat_id, m.sender_id, m.content, COALESCE(EXTRACT(EPOCH FROM m.sent_at)::bigint * 1000, 0) AS sent_at, m.is_delivered
		FROM messages m
		JOIN chat_participants cp ON m.chat_id = cp.chat_id
		WHERE cp.user_id = $1
		ORDER BY m.sent_at DESC
	`

	err := r.db.Select(&messages, query, userID)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("failed to get user messages: %v", err))
	}
	result := make([]*models.Message, len(messages))
	for i := range messages {
		result[i] = &messages[i]
	}
	return result, nil
}

func (r *MessageRepo) GetMessagesByChatID(chatID string) ([]*models.Message, error) {
	var messages []models.Message
	query := "SELECT id, chat_id, sender_id, content, COALESCE(EXTRACT(EPOCH FROM sent_at)::bigint * 1000, 0) AS sent_at, is_delivered FROM messages WHERE chat_id = $1 ORDER BY sent_at ASC"
	err := r.db.Select(&messages, query, chatID)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("failed to get user messages: %v", err))
	}
	result := make([]*models.Message, len(messages))
	for i := range messages {
		result[i] = &messages[i]
	}
	return result, nil
}

func (r *MessageRepo) UsersSendMess(chatID string, senderID string) (*[]string, error) {
	var users []string
	query := "SELECT user_id FROM chat_participants WHERE chat_id = $1 AND user_id != $2"
	err := r.db.Select(&users, query, chatID, senderID)
	if err != nil {
		return nil, errors.NewCustomError(http.StatusInternalServerError, fmt.Sprintf("failed to get user messages: %v", err))
	}
	return &users, nil
}
