package models

type WebSocketRequest struct {
	Type    string `json:"type"`
	Message byte   `json:"message"`
}
