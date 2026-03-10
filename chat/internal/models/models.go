package models

type CreatePrivateChatRequest struct {
	FriendID string `json:"friend_id" required:"true"`
}
type CreatePrivateChatResponse struct {
	ChatID string `json:"chat_id"`
}

type CreatePublicChatRequest struct {
	Name           string   `json:"name" required:"true"`
	ParticipantIDs []string `json:"participant_id" required:"true"`
}
type CreatePublicChatResponse struct {
	ChatID string `json:"chat_id"`
}

type GetChatsRequest struct{}
type GetChatsResponse struct {
	ChatID []string `json:"chat_id"`
}

type GetChatUsersRequest struct{}
type GetChatUsersResponse struct {
	Users []string `json:"users"`
}
