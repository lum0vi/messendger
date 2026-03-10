import { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { chatApi, userApi } from '@/api';
import { useAuth } from '@/context/AuthContext';
import { Button, Card, Input } from '@/components/ui';
import type { User } from '@/types';

export function ChatsPage() {
  const navigate = useNavigate();
  const { logout } = useAuth();
  const [chatIds, setChatIds] = useState<string[]>([]);
  const [users, setUsers] = useState<User[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const [showNewChatModal, setShowNewChatModal] = useState(false);
  const [newChatType, setNewChatType] = useState<'private' | 'public'>('private');
  const [selectedUserId, setSelectedUserId] = useState('');
  const [publicChatName, setPublicChatName] = useState('');
  const [selectedParticipants, setSelectedParticipants] = useState<string[]>([]);
  const [isCreating, setIsCreating] = useState(false);

  useEffect(() => {
    loadData();
  }, []);

  const loadData = async () => {
    try {
      const [chatsRes, usersRes] = await Promise.all([
        chatApi.getMyChats(),
        userApi.getUsers(),
      ]);
      setChatIds(chatsRes.chat_id || []);
      setUsers(usersRes.Users);
    } catch (error) {
      console.error('Failed to load data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const handleCreateChat = async () => {
    setIsCreating(true);
    try {
      let response;
      if (newChatType === 'private') {
        if (!selectedUserId) {
          alert('Выберите пользователя');
          setIsCreating(false);
          return;
        }
        response = await chatApi.createPrivateChat({ friend_id: selectedUserId });
      } else {
        if (!publicChatName) {
          alert('Введите название чата');
          setIsCreating(false);
          return;
        }
        response = await chatApi.createPublicChat({
          name: publicChatName,
          participant_id: selectedParticipants,
        });
      }
      
      navigate(`/chat/${response.chat_id}`);
    } catch (error) {
      console.error('Failed to create chat:', error);
      alert('Ошибка при создании чата');
    } finally {
      setIsCreating(false);
    }
  };

  const toggleParticipant = (userId: string) => {
    setSelectedParticipants(prev =>
      prev.includes(userId)
        ? prev.filter(id => id !== userId)
        : [...prev, userId]
    );
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
            <Button variant="secondary" onClick={() => navigate('/profile')}>
              Профиль
            </Button>
            <Button variant="danger" onClick={logout}>
              Выйти
            </Button>
          </div>
        </div>
      </header>

      <main className="max-w-7xl mx-auto px-4 py-6">
        <div className="flex justify-between items-center mb-6">
          <h2 className="text-2xl font-bold">Чаты</h2>
          <Button onClick={() => setShowNewChatModal(true)}>
            Новый чат
          </Button>
        </div>

        {chatIds.length === 0 ? (
          <Card className="text-center py-8">
            <p className="text-gray-500 mb-4">У вас пока нет чатов</p>
            <Button onClick={() => setShowNewChatModal(true)}>
              Создать первый чат
            </Button>
          </Card>
        ) : (
          <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
            {chatIds.map((chatId) => (
              <Card
                key={chatId}
                className="cursor-pointer hover:shadow-lg transition-shadow"
                onClick={() => navigate(`/chat/${chatId}`)}
              >
                <div className="flex items-center gap-3">
                  <div className="w-12 h-12 bg-blue-500 rounded-full flex items-center justify-center text-white font-bold">
                    {chatId.substring(0, 2).toUpperCase()}
                  </div>
                  <div>
                    <p className="font-medium">Чат</p>
                    <p className="text-sm text-gray-500">{chatId.substring(0, 8)}...</p>
                  </div>
                </div>
              </Card>
            ))}
          </div>
        )}
      </main>

      {/* New Chat Modal */}
      {showNewChatModal && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <Card className="w-full max-w-md mx-4">
            <h3 className="text-xl font-bold mb-4">Создать чат</h3>
            
            <div className="flex gap-2 mb-4">
              <Button
                variant={newChatType === 'private' ? 'primary' : 'secondary'}
                onClick={() => setNewChatType('private')}
                className="flex-1"
              >
                Личный
              </Button>
              <Button
                variant={newChatType === 'public' ? 'primary' : 'secondary'}
                onClick={() => setNewChatType('public')}
                className="flex-1"
              >
                Групповой
              </Button>
            </div>

            {newChatType === 'private' ? (
              <div className="mb-4">
                <label className="block text-sm font-medium text-gray-700 mb-1">
                  Выберите пользователя
                </label>
                <select
                  className="w-full px-3 py-2 border border-gray-300 rounded-md"
                  value={selectedUserId}
                  onChange={(e) => setSelectedUserId(e.target.value)}
                >
                  <option value="">-- Выберите --</option>
                  {users.map((user) => (
                    <option key={user.id} value={user.id}>
                      {user.username} ({user.email})
                    </option>
                  ))}
                </select>
              </div>
            ) : (
              <>
                <Input
                  label="Название чата"
                  type="text"
                  value={publicChatName}
                  onChange={(e) => setPublicChatName(e.target.value)}
                  placeholder="Введите название"
                />
                <div className="mb-4">
                  <label className="block text-sm font-medium text-gray-700 mb-1">
                    Участники
                  </label>
                  <div className="max-h-40 overflow-y-auto border border-gray-300 rounded-md p-2">
                    {users.map((user) => (
                      <label
                        key={user.id}
                        className="flex items-center gap-2 py-1 cursor-pointer"
                      >
                        <input
                          type="checkbox"
                          checked={selectedParticipants.includes(user.id)}
                          onChange={() => toggleParticipant(user.id)}
                          className="rounded"
                        />
                        <span>{user.username}</span>
                      </label>
                    ))}
                  </div>
                </div>
              </>
            )}

            <div className="flex gap-2">
              <Button
                variant="secondary"
                onClick={() => setShowNewChatModal(false)}
                className="flex-1"
              >
                Отмена
              </Button>
              <Button
                onClick={handleCreateChat}
                isLoading={isCreating}
                className="flex-1"
              >
                Создать
              </Button>
            </div>
          </Card>
        </div>
      )}
    </div>
  );
}
