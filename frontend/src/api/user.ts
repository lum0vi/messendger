import api from './client';
import type {
  User,
  UpdateUserRequest,
  GetUserByIdRequest,
  GetUserByUsernameRequest,
} from '@/types';

export const userApi = {
  getMe: async (): Promise<User> => {
    const response = await api.get<User>('/user/me');
    return response.data;
  },

  updateMe: async (data: UpdateUserRequest): Promise<void> => {
    await api.put('/user/me', data);
  },

  getUsers: async (): Promise<{ Users: User[] }> => {
    const response = await api.get<{ Users: User[] }>('/user/users');
    return response.data;
  },

  getUserById: async (data: GetUserByIdRequest): Promise<User> => {
    const response = await api.post<User>('/user/id', data);
    return response.data;
  },

  getUserByUsername: async (data: GetUserByUsernameRequest): Promise<User> => {
    const response = await api.post<User>('/user/name', data);
    return response.data;
  },
};
