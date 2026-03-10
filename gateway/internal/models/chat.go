package models

type CreatePrivateChatRequest struct {
	FriendId string `json:"friend_id"`
}
type CreatePrivateChatResponse struct {
	ChatId string `json:"chat_id"`
}

type CreatePublicChatRequest struct {
	Name          string   `json:"name"`
	ParticipantId []string `json:"participant_id"`
}
type CreatePublicChatResponse struct {
	ChatId string `json:"chat_id"`
}

type GetMeChatsRequest struct{}
type GetMeChatsResponse struct {
	ChatId []string `json:"chat_id"`
}

type GetChatUsersRequest struct{}
type GetChatUsersResponse struct {
	Users []string `json:"users"`
}
