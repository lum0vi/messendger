import api from './client';
import type { Message } from '@/types';

export const messageApi = {
  updateMessageStatus: async (): Promise<void> => {
    await api.put('/message/upd', {});
  },

  getUnsentMessages: async (messageId: string): Promise<void> => {
    await api.get(`/message/${messageId}`);
  },

  getChatMessages: async (chatId: string): Promise<Message[]> => {
    const response = await api.get(`/message/chat/${chatId}`);
    return response.data.messages || [];
  },
};
