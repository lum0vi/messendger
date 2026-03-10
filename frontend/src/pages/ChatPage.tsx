import React, { useState, useEffect, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { chatApi, messageApi } from '@/api';
import { useAuth } from '@/context/AuthContext';
import { useWebSocket } from '@/hooks/useWebSocket';
import { Button, Card } from '@/components/ui';
import type { Message, WebSocketMessage } from '@/types';

export function ChatPage() {
  const { chatId } = useParams<{ chatId: string }>();
  const navigate = useNavigate();
  const { user } = useAuth();
  const [messages, setMessages] = useState<Message[]>([]);
  const [newMessage, setNewMessage] = useState('');
  const [participants, setParticipants] = useState<string[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const messagesEndRef = useRef<HTMLDivElement>(null);

  const handleWebSocketMessage = (message: WebSocketMessage) => {
    if (message.type === 'new_message') {
      const msg = message.data as Message;
      if (msg.chat_id === chatId) {
        setMessages((prev) => [...prev, msg]);
      }
    }
  };

  const { isConnected, connect, sendMessage: sendWsMessage } = useWebSocket({
    onMessage: handleWebSocketMessage,
  });

  useEffect(() => {
    if (chatId) {
      loadChatData();
      connect();
    }
  }, [chatId]);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const loadChatData = async () => {
    try {
      const [chatUsersResponse, messagesResponse] = await Promise.all([
        chatApi.getChatUsers(chatId!),
        messageApi.getChatMessages(chatId!),
      ]);
      setParticipants(chatUsersResponse.users);
      setMessages(messagesResponse);
    } catch (error) {
      console.error('Failed to load chat data:', error);
    } finally {
      setIsLoading(false);
    }
  };

  const scrollToBottom = () => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' });
  };

  const handleSendMessage = (e: React.FormEvent) => {
    e.preventDefault();
    if (!newMessage.trim() || !chatId || !user) return;

    const message: Message = {
      chat_id: chatId,
      sender_id: user.id,
      content: newMessage.trim(),
      sent_at: Date.now(),
    };

    sendWsMessage({
      type: 'send_message',
      data: message,
    });

    setMessages((prev) => [...prev, message]);
    setNewMessage('');
  };

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleTimeString('ru-RU', {
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
    <div className="min-h-screen bg-gray-100 flex flex-col">
      {/* Header */}
      <header className="bg-white shadow">
        <div className="max-w-7xl mx-auto px-4 py-4 flex items-center gap-4">
          <Button variant="secondary" onClick={() => navigate('/chats')}>
            ← Назад
          </Button>
          <div className="flex-1">
            <h1 className="text-lg font-bold">Чат</h1>
            <p className="text-sm text-gray-500">
              {participants.length} участников
              {isConnected && (
                <span className="ml-2 text-green-500">● Подключено</span>
              )}
            </p>
          </div>
        </div>
      </header>

      {/* Messages */}
      <main className="flex-1 overflow-y-auto p-4">
        <div className="max-w-3xl mx-auto space-y-4">
          {messages.length === 0 ? (
            <Card className="text-center py-8">
              <p className="text-gray-500">Нет сообщений. Начните диалог!</p>
            </Card>
          ) : (
            messages.map((msg, index) => {
              const isOwn = msg.sender_id === user?.id;
              return (
                <div
                  key={index}
                  className={`flex ${isOwn ? 'justify-end' : 'justify-start'}`}
                >
                  <div
                    className={`max-w-xs lg:max-w-md px-4 py-2 rounded-lg ${
                      isOwn
                        ? 'bg-blue-500 text-white'
                        : 'bg-white shadow'
                    }`}
                  >
                    {!isOwn && (
                      <p className="text-xs text-gray-500 mb-1">
                        {msg.sender_id.substring(0, 8)}...
                      </p>
                    )}
                    <p>{msg.content}</p>
                    <p
                      className={`text-xs mt-1 ${
                        isOwn ? 'text-blue-100' : 'text-gray-400'
                      }`}
                    >
                      {msg.sent_at && formatTime(msg.sent_at)}
                    </p>
                  </div>
                </div>
              );
            })
          )}
          <div ref={messagesEndRef} />
        </div>
      </main>

      {/* Message Input */}
      <footer className="bg-white border-t">
        <form
          onSubmit={handleSendMessage}
          className="max-w-3xl mx-auto px-4 py-4 flex gap-2"
        >
          <input
            type="text"
            value={newMessage}
            onChange={(e) => setNewMessage(e.target.value)}
            placeholder="Введите сообщение..."
            className="flex-1 px-4 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
          <Button type="submit" disabled={!newMessage.trim()}>
            Отправить
          </Button>
        </form>
      </footer>
    </div>
  );
}
