import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { userApi } from '@/api';
import { useAuth } from '@/context/AuthContext';
import { Button, Card, Input } from '@/components/ui';
import type { User, UpdateUserRequest } from '@/types';

export function ProfilePage() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [user, setUser] = useState<User | null>(null);
  const [isLoading, setIsLoading] = useState(true);
  const [isEditing, setIsEditing] = useState(false);
  const [isSaving, setIsSaving] = useState(false);
  const [formData, setFormData] = useState<UpdateUserRequest>({
    username: '',
    email: '',
    password: '',
  });
  const [error, setError] = useState('');
  const [success, setSuccess] = useState('');

  useEffect(() => {
    loadUser();
  }, []);

  const loadUser = async () => {
    try {
      const userData = await userApi.getMe();
      setUser(userData);
      setFormData({
        username: userData.username,
        email: userData.email,
        password: '',
      });
    } catch (error) {
      console.error('Failed to load user:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
    setError('');
    setSuccess('');
  };

  const handleSave = async () => {
    setIsSaving(true);
    setError('');
    setSuccess('');

    try {
      const updateData: UpdateUserRequest = {};
      if (formData.username && formData.username !== user?.username) {
        updateData.username = formData.username;
      }
      if (formData.email && formData.email !== user?.email) {
        updateData.email = formData.email;
      }
      if (formData.password) {
        updateData.password = formData.password;
      }

      if (Object.keys(updateData).length > 0) {
        await userApi.updateMe(updateData);
        setSuccess('Профиль успешно обновлен');
        setIsEditing(false);
        await loadUser();
      }
    } catch (error) {
      setError('Ошибка при обновлении профиля');
    } finally {
      setIsSaving(false);
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('ru-RU', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (isLoading) {
    return (
      <div className="min-h-screen flex items-center justify-center">
        <div className="text-lg">Загрузка...</div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-100">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex justify-between items-center">
          <h1 className="text-xl font-bold">Messenger</h1>
          <div className="flex gap-2">
            <Button variant="secondary" onClick={() => navigate('/chats')}>
              ← К чатам
            </Button>
            <Button variant="danger" onClick={logout}>
              Выйти
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-2xl mx-auto px-4 py-6">
        <Card>
          <div className="flex justify-between items-center mb-6">
            <h2 className="text-2xl font-bold">Профиль</h2>
            {!isEditing && (
              <Button onClick={() => setIsEditing(true)}>
                Редактировать
              </Button>
            )}
          </div>

          {error && (
            <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md text-sm">
              {error}
            </div>
          )}

          {success && (
            <div className="mb-4 p-3 bg-green-100 text-green-700 rounded-md text-sm">
              {success}
            </div>
          )}

          {isEditing ? (
            <div>
              <Input
                label="Имя пользователя"
                type="text"
                name="username"
                value={formData.username}
                onChange={handleChange}
              />
              <Input
                label="Email"
                type="email"
                name="email"
                value={formData.email}
                onChange={handleChange}
              />
              <Input
                label="Новый пароль (оставьте пустым, если не хотите менять)"
                type="password"
                name="password"
                value={formData.password}
                onChange={handleChange}
              />
              <div className="flex gap-2 mt-4">
                <Button
                  variant="secondary"
                  onClick={() => {
                    setIsEditing(false);
                    setFormData({
                      username: user?.username || '',
                      email: user?.email || '',
                      password: '',
                    });
                    setError('');
                  }}
                  className="flex-1"
                >
                  Отмена
                </Button>
                <Button
                  onClick={handleSave}
                  isLoading={isSaving}
                  className="flex-1"
                >
                  Сохранить
                </Button>
              </div>
            </div>
          ) : (
            <div className="space-y-4">
              <div>
                <p className="text-sm text-gray-500">ID</p>
                <p className="font-mono text-sm">{user?.id}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Имя пользователя</p>
                <p className="text-lg">{user?.username}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Email</p>
                <p className="text-lg">{user?.email}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Дата регистрации</p>
                <p className="text-lg">{user?.created_at && formatDate(user.created_at)}</p>
              </div>
              <div>
                <p className="text-sm text-gray-500">Последнее обновление</p>
                <p className="text-lg">{user?.updated_at && formatDate(user.updated_at)}</p>
              </div>
            </div>
          )}
        </Card>

        {/* Users List */}
        <Card className="mt-6">
          <h3 className="text-xl font-bold mb-4">Пользователи</h3>
          <UsersList />
        </Card>
      </main>
    </div>
  );
}

function UsersList() {
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [searchId, setSearchId] = useState('');
  const [searchUsername, setSearchUsername] = useState('');
  const [searchResult, setSearchResult] = useState<User | null>(null);
  const [searchError, setSearchError] = useState('');

  useEffect(() => {
    loadUsers();
  }, []);

  const loadUsers = async () => {
    try {
      const response = await userApi.getUsers();
      setUsers(response.Users);
    } catch (error) {
      console.error('Failed to load users:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleSearchById = async () => {
    if (!searchId.trim()) return;
    setSearchError('');
    setSearchResult(null);
    try {
      const result = await userApi.getUserById({ id: searchId });
      setSearchResult(result);
    } catch {
      setSearchError('Пользователь не найден');
    }
  };

  const handleSearchByUsername = async () => {
    if (!searchUsername.trim()) return;
    setSearchError('');
    setSearchResult(null);
    try {
      const result = await userApi.getUserByUsername({ username: searchUsername });
      setSearchResult(result);
    } catch {
      setSearchError('Пользователь не найден');
    }
  };

  if (isLoading) {
    return <div>Загрузка пользователей...</div>;
  }

  return (
    <div>
      {/* Search by ID */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Поиск по ID
        </label>
        <div className="flex gap-2">
          <input
            type="text"
            value={searchId}
            onChange={(e) => setSearchId(e.target.value)}
            placeholder="Введите ID пользователя"
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md"
          />
          <Button onClick={handleSearchById}>Найти</Button>
        </div>
      </div>

      {/* Search by Username */}
      <div className="mb-4">
        <label className="block text-sm font-medium text-gray-700 mb-1">
          Поиск по имени пользователя
        </label>
        <div className="flex gap-2">
          <input
            type="text"
            value={searchUsername}
            onChange={(e) => setSearchUsername(e.target.value)}
            placeholder="Введите имя пользователя"
            className="flex-1 px-3 py-2 border border-gray-300 rounded-md"
          />
          <Button onClick={handleSearchByUsername}>Найти</Button>
        </div>
      </div>

      {searchError && (
        <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-md text-sm">
          {searchError}
        </div>
      )}

      {searchResult && (
        <div className="mb-4 p-3 bg-blue-50 rounded-md">
          <p className="font-medium">{searchResult.username}</p>
          <p className="text-sm text-gray-500">{searchResult.email}</p>
        </div>
      )}

      <h4 className="font-medium mb-2">Все пользователи</h4>
      <div className="max-h-60 overflow-y-auto space-y-2">
        {users.map((u) => (
          <div
            key={u.id}
            className="p-3 bg-gray-50 rounded-md flex justify-between items-center"
          >
            <div>
              <p className="font-medium">{u.username}</p>
              <p className="text-sm text-gray-500">{u.email}</p>
            </div>
            <span className="text-xs text-gray-400 font-mono">
              {u.id.substring(0, 8)}...
            </span>
          </div>
        ))}
      </div>
    </div>
  );
}
