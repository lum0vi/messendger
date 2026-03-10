// Auth types
export interface RegisterRequest {
  username: string;
  password: string;
  email: string;
}

export interface RegisterResponse {
  id: string;
}

export interface LoginRequest {
  username: string;
  password: string;
}

export interface LoginResponse {
  token: string;
}

// User types
export interface User {
  id: string;
  username: string;
  email: string;
  created_at: string;
  updated_at: string;
}

export interface UpdateUserRequest {
  username?: string;
  password?: string;
  email?: string;
}

export interface GetUserByIdRequest {
  id: string;
}

export interface GetUserByUsernameRequest {
  username: string;
}

// Chat types
export interface CreatePrivateChatRequest {
  friend_id: string;
}

export interface CreatePublicChatRequest {
  name: string;
  participant_id: string[];
}

export interface ChatResponse {
  chat_id: string;
}

export interface ChatsResponse {
  chat_id: string[] | null;
}

export interface ChatUsersResponse {
  users: string[];
}

// Message types
export interface Message {
  id?: string;
  chat_id: string;
  sender_id: string;
  content: string;
  sent_at?: number;
}

export interface WebSocketMessage {
  type: 'new_message' | 'user_joined' | 'user_left' | 'message_delivered';
  data: Message | unknown;
}

// API Error
export interface ApiError {
  code: number;
  message: string;
}
