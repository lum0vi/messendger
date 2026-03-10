import api from './client';
import type {
  CreatePrivateChatRequest,
  CreatePublicChatRequest,
  ChatResponse,
  ChatsResponse,
  ChatUsersResponse,
} from '@/types';

export const chatApi = {
  getMyChats: async (): Promise<ChatsResponse> => {
    const response = await api.get<ChatsResponse>('/chat/');
    return response.data;
  },

  createPrivateChat: async (data: CreatePrivateChatRequest): Promise<ChatResponse> => {
    const response = await api.post<ChatResponse>('/chat/private', data);
    return response.data;
  },

  createPublicChat: async (data: CreatePublicChatRequest): Promise<ChatResponse> => {
    const response = await api.post<ChatResponse>('/chat/public', data);
    return response.data;
  },

  getChatUsers: async (chatId: string): Promise<ChatUsersResponse> => {
    const response = await api.get<ChatUsersResponse>(`/chat/${chatId}/users`);
    return response.data;
  },
};
