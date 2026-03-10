package models

import "time"

type RegisterRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required"`
}

type RegisterResponse struct {
	ID string `json:"id"`
}

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type GetMeRequest struct {
}
type GetMeResponse struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type UpdateMeRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}
type UpdateMeResponse struct{}

type GetUsersRequest struct {
}

type GetUsersResponse struct {
	Users []*UserForGetUsers
}

type UserForGetUsers struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type GetUserByIdRequest struct {
	ID string `json:"id"`
}
type GetUserByIdResponse struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type GetUserByUsernameRequest struct {
	Username string `json:"username"`
}
type GetUserByUsernameResponse struct {
	ID       string    `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Password string    `json:"password"`
	CreateAt time.Time `json:"created_at"`
	UpdateAt time.Time `json:"updated_at"`
}

type WSInput struct {
}

type GetUserMessageResponse struct {
	Messages []*Message `json:"messages"`
}
