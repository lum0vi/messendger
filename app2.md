# Детальное описание работы сервисов Messenger Application

## Содержание

1. [Auth Service (Порт 8081)](#auth-service-порт-8081)
2. [User Service (Порт 8082)](#user-service-порт-8082)
3. [Chat Service (Порт 8083)](#chat-service-порт-8083)
4. [Message Service (Порт 8084)](#message-service-порт-8084)
5. [Gateway (Порт 8080)](#gateway-порт-8080)
6. [Frontend (Порт 3000)](#frontend-порт-3000)

---

## Auth Service (Порт 8081)

Сервис отвечает за аутентификацию и регистрацию пользователей.

### Эндпоинты

#### POST /auth/register — Регистрация пользователя

**Входящие данные (Request Body):**
```json
{
    "username": "string",  // Обязательно, уникальное
    "password": "string",  // Обязательно
    "email": "string"      // Обязательно, уникальное
}
```

**Процесс выполнения:**

1. **Handler** (`auth/internal/handler/auth.go:Register`):
   - Извлекает `request_id` из контекста
   - Десериализует JSON в `models.RegisterRequest`
   - Вызывает `h.srv.Auth.Register(&req)`

2. **Service** (`auth/internal/service/auth.go:Register`):
   - Хэширует пароль: `hash.HashPass(req.Password)` (использует bcrypt)
   - Передаёт данные в Repository: `s.repo.PostgresAuth.Register(req)`
   - Возвращает `id` созданного пользователя

3. **Repository** (`auth/internal/repository/postgres_auth.go:Register`):
   - Выполняет SQL запрос:
     ```sql
     INSERT INTO users (username, password_hash, email) VALUES ($1, $2, $3) RETURNING id
     ```
   - При ошибке уникальности username возвращает HTTP 409 Conflict

**Ответ (Response):**
```json
{
    "id": "uuid-пользователя"
}
```

**Ошибки:**
- 400 — некорректные данные (пустые поля)
- 409 — username уже занят
- 500 — внутренняя ошибка сервера

---

#### POST /auth/login — Вход пользователя

**Входящие данные (Request Body):**
```json
{
    "username": "string",  // Обязательно
    "password": "string"   // Обязательно
}
```

**Процесс выполнения:**

1. **Handler** (`auth/internal/handler/auth.go:Login`):
   - Десериализует JSON в `models.LoginRequest`
   - Вызывает `h.srv.Auth.Login(&req)`

2. **Service** (`auth/internal/service/auth.go:Login`):
   - Получает `id` пользователя: `s.repo.PostgresAuth.Login(req)`
   - Создаёт JWT токен: `jwt.CreateJWT(id)`
   - Возвращает токен

3. **Repository** (`auth/internal/repository/postgres_auth.go:Login`):
   - Выполняет SQL запрос:
     ```sql
     SELECT id, password_hash FROM users WHERE username = $1
     ```
   - Если пользователь не найден — ошибка
   - Проверяет пароль: `hash.VerifyPass(req.Password, passHash)` (bcrypt compare)
   - При неверном пароле возвращает HTTP 401 Unauthorized

4. **JWT Создание** (`auth/internal/jwt/jwt.go:CreateJWT`):
   - Создаёт claims с `user_id` и `exp` (24 часа)
   - Подписывает токен приватным RSA ключом (RS256)

**Ответ (Response):**
```json
{
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Ошибки:**
- 400 — некорректные данные
- 401 — неверный пароль или пользователь не найден
- 500 — внутренняя ошибка сервера

---

## User Service (Порт 8082)

Сервис управления пользователями. Все endpoints требуют заголовок `id` с ID пользователя (из JWT токена).

### Эндпоинты

#### GET /user/me — Получить данные текущего пользователя

**Входящие заголовки:**
```
id: uuid-пользователя (из JWT)
Authorization: Bearer <token>
```

**Процесс выполнения:**

1. **Handler** (`user/internal/handler/user.go:GetMe`):
   - Извлекает `id` из заголовка
   - Вызывает `h.svc.User.GetMe(id)`

2. **Service**:
   - Вызывает repository: `r.User.GetMe(id)`

3. **Repository**:
   - Выполняет SQL запрос:
     ```sql
     SELECT id, username, email, password_hash, created_at, updated_at 
     FROM users WHERE id = $1
     ```

**Ответ (Response):**
```json
{
    "id": "uuid",
    "username": "string",
    "email": "string",
    "password_hash": "bcrypt-hash",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

---

#### PUT /user/me — Обновить профиль пользователя

**Входящие заголовки:**
```
id: uuid-пользователя
Content-Type: application/json
```

**Входящие данные (Request Body):**
```json
{
    "username": "string",  // Опционально
    "email": "string",     // Опционально
    "password": "string"   // Опционально
}
```

**Процесс выполнения:**

1. **Handler** (`user/internal/handler/user.go:UpdateMe`):
   - Извлекает `id` из заголовка
   - Десериализует JSON в `models.UpdateMeRequest`
   - Вызывает `h.svc.User.UpdateMe(id, &req)`

2. **Service**:
   - Если передан новый пароль — хэширует его
   - Вызывает repository для обновления

3. **Repository**:
   - Выполняет SQL запрос (формируется динамически):
     ```sql
     UPDATE users SET 
       username = COALESCE($2, username),
       email = COALESCE($3, email),
       password_hash = COALESCE($4, password_hash),
       updated_at = NOW()
     WHERE id = $1
     ```

**Ответ (Response):**
```json
{}
```

---

#### GET /user/users — Получить всех пользователей

**Входящие заголовки:**
```
id: uuid-пользователя
Authorization: Bearer <token>
```

**Процесс выполнения:**

1. **Handler** (`user/internal/handler/user.go:GetUsers`):
   - Вызывает `h.svc.User.GetUsers()`

2. **Repository**:
   - Выполняет SQL запрос:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users
     ```

**Ответ (Response):**
```json
{
    "Users": [
        {
            "id": "uuid",
            "username": "string",
            "email": "string",
            "created_at": "timestamp",
            "updated_at": "timestamp"
        },
        ...
    ]
}
```

---

#### POST /user/id — Получить пользователя по ID

**Входящие заголовки:**
```
id: uuid-запрашивающего пользователя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные (Request Body):**
```json
{
    "id": "uuid-пользователя"
}
```

**Процесс выполнения:**

1. **Handler** (`user/internal/handler/user.go:GetUserById`):
   - Десериализует JSON в `models.GetUserByIDRequest`
   - Вызывает `h.svc.User.GetUserByID(&req)`

2. **Repository**:
   - Выполняет SQL запрос:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users WHERE id = $1
     ```

**Ответ (Response):**
```json
{
    "id": "uuid",
    "username": "string",
    "email": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

---

#### POST /user/name — Получить пользователя по имени

**Входящие данные (Request Body):**
```json
{
    "username": "string"
}
```

**Процесс выполнения:**

1. **Handler** (`user/internal/handler/user.go:GetUserByName`):
   - Десериализует JSON в `models.GetUserByUsernameRequest`
   - Вызывает `h.svc.User.GetUserByUsername(&req)`

2. **Repository**:
   - Выполняет SQL запрос:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users WHERE username = $1
     ```

**Ответ (Response):**
```json
{
    "id": "uuid",
    "username": "string",
    "email": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

---

## Chat Service (Порт 8083)

Сервис управления чатами. Требует заголовок `id` с ID пользователя.

### Эндпоинты

#### POST /chat/private — Создать приватный чат

**Входящие данные (Request Body):**
```json
{
    "friend_id": "uuid-другого-пользователя"
}
```

**Процесс выполнения:**

1. **Handler** (`chat/internal/handler/chat.go:CreatePrivateChat`):
   - Извлекает `id` (создателя) из заголовка
   - Десериализует JSON в `models.CreatePrivateChatRequest`
   - Вызывает `h.svc.Chat.CreatePrivateChat(id, req)`

2. **Service**:
   - Вызывает repository

3. **Repository** (`chat/internal/repository/chat.go:CreatePrivateChat`):
   - Начинает транзакцию
   - Создаёт запись в таблице `chats`:
     ```sql
     INSERT INTO chats (is_private) VALUES ($1) RETURNING id
     ```
     - `is_private = true`
   - Добавляет создателя в участники:
     ```sql
     INSERT INTO chat_participants (chat_id, user_id) VALUES ($1, $2)
     ```
   - Добавляет друга в участники:
     ```sql
     INSERT INTO chat_participants (chat_id, user_id) VALUES ($1, $2)
     ```
   - Коммитит транзакцию

**Ответ (Response):**
```json
{
    "chat_id": "uuid-чата"
}
```

**Таблицы БД:**
- `chats`: создаётся 1 запись
- `chat_participants`: создаётся 2 записи

---

#### POST /chat/public — Создать публичный чат

**Входящие данные (Request Body):**
```json
{
    "name": "Название чата",
    "participant_id": ["uuid1", "uuid2", ...]
}
```

**Процесс выполнения:**

1. **Handler** (`chat/internal/handler/chat.go:CreatePublicChat`):
   - Извлекает `id` создателя из заголовка
   - Десериализует JSON в `models.CreatePublicChatRequest`
   - Вызывает `h.svc.Chat.CreatePublicChat(id, req)`

2. **Repository** (`chat/internal/repository/chat.go:CreatePublicChat`):
   - Начинает транзакцию
   - Создаёт запись в таблице `chats`:
     ```sql
     INSERT INTO chats (is_private, name) VALUES ($1, $2) RETURNING id
     ```
     - `is_private = false`
     - `name = req.Name`
   - Добавляет всех участников (включая создателя):
     ```sql
     INSERT INTO chat_participants (chat_id, user_id) VALUES ($1, $2)
     ```
   - Коммитит транзакцию

**Ответ (Response):**
```json
{
    "chat_id": "uuid-чата"
}
```

**Таблицы БД:**
- `chats`: создаётся 1 запись
- `chat_participants`: создаётся N записей (количество участников)

---

#### GET /chat/ — Получить чаты пользователя

**Процесс выполнения:**

1. **Handler** (`chat/internal/handler/chat.go:GetChats`):
   - Извлекает `id` пользователя из заголовка
   - Вызывает `h.svc.Chat.GetChats(id)`

2. **Repository** (`chat/internal/repository/chat.go:GetChats`):
   - Выполняет SQL запрос:
     ```sql
     SELECT chat_id FROM chat_participants WHERE user_id = $1
     ```

**Ответ (Response):**
```json
{
    "chat_id": ["uuid-чата1", "uuid-чата2", ...]
}
```

---

#### GET /chat/:chat_id/users — Получить участников чата

**Параметры URL:**
- `chat_id` — UUID чата

**Процесс выполнения:**

1. **Handler** (`chat/internal/handler/chat.go:GetChatUsers`):
   - Извлекает `chat_id` из параметров URL
   - Вызывает `h.svc.Chat.GetUsersChat(chatID)`

2. **Repository** (`chat/internal/repository/chat.go:GetUsersChat`):
   - Выполняет SQL запрос:
     ```sql
     SELECT user_id FROM chat_participants WHERE chat_id = $1
     ```

**Ответ (Response):**
```json
{
    "users": ["uuid-user1", "uuid-user2", ...]
}
```

---

## Message Service (Порт 8084)

Сервис обработки сообщений. Работает полностью асинхронно через Apache Kafka.

### Kafka Топики

#### Входящий: `message.send`
Сообщения отправляются из Gateway при отправке пользователем через WebSocket.

#### Исходящий: `message.delivered`
Уведомления о доставке отправляются в Gateway для пересылки получателю.

### Процесс обработки сообщений

#### 1. Чтение из Kafka

**Consumer** (`message/internal/kafka/consumer.go:Start`):
- Подписывается на топик `message.send`
- Читает сообщения в бесконечном цикле

**Входящее сообщение из Kafka:**
```json
{
    "type": "send_message",
    "data": {
        "chat_id": "uuid-чата",
        "sender_id": "uuid-отправителя",
        "content": "текст сообщения",
        "sent_at": 1234567890
    }
}
```

#### 2. Сохранение в БД

**Repository** (`message/internal/repository/message.go:Save`):
- Выполняет SQL запрос:
  ```sql
  INSERT INTO messages (chat_id, sender_id, content) 
  VALUES ($1, $2, $3) 
  RETURNING id, EXTRACT(EPOCH FROM sent_at)::bigint * 1000
  ```
- Возвращает `id` сообщения и `sent_at` timestamp

#### 3. Определение получателей

**Repository** (`message/internal/repository/message.go:UsersSendMess`):
- Выполняет SQL запрос:
  ```sql
  SELECT user_id FROM chat_participants 
  WHERE chat_id = $1 AND user_id != $2
  ```
- Возвращает всех участников чата, кроме отправителя

#### 4. Отправка уведомлений

**Producer** (`message/internal/kafka/producer.go`):
- Для каждого получателя отправляет сообщение в топик `message.delivered`

**Сообщение в топике `message.delivered`:**
```json
{
    "user_id": "uuid-получателя",
    "chat_id": "uuid-чата",
    "sender_id": "uuid-отправителя",
    "content": "текст сообщения",
    "sent_at": 1234567890
}
```

### Таблицы БД

**messages:**
| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | PK, автогенерируется |
| chat_id | UUID | FK на chats |
| sender_id | UUID | FK на users |
| content | TEXT | Текст сообщения |
| is_delivered | BOOLEAN | Статус доставки |
| sent_at | TIMESTAMPTZ | Время отправки |

---

## Gateway (Порт 8080)

API Gateway — центральная точка входа. Объединяет все сервисы, проверяет JWT, управляет WebSocket соединениями.

### Эндпоинты

#### GET /ws — WebSocket соединение

**Запрос:**
```
GET ws://localhost:8080/ws?token=<JWT>
```

**Процесс выполнения:**

1. **Middleware** (`gateway/internal/middleware/auth.go:AuthMiddleware`):
   - Извлекает токен из query параметра `token`
   - Проверяет JWT токен
   - Извлекает `user_id` из claims
   - Сохраняет `user_id` в контексте gin

2. **Handler** (`gateway/internal/handler/websocket.go:Websocket`):
   - Создаёт `socketID` (UUID)
   - Сохраняет в Redis:
     - `socket:<user_id>` -> `socketID` (TTL 30 минут)
     - `socket_to_user:<socketID>` -> `user_id` (TTL 30 минут)
   - Обновляет HTTP соединение до WebSocket
   - Добавляет соединение в мапу `h.connections[socketID]`

3. **Обработка сообщений**:
   - Читает сообщение от клиента
   - Пересылает в Kafka топик `message.send`:
     ```json
     {
         "type": "send_message",
         "data": {
             "chat_id": "uuid-чата",
             "sender_id": "uuid-отправителя",
             "content": "текст сообщения",
             "sent_at": 1234567890
         }
     }
     ```

4. **Получение сообщений**:
   - **Consumer** (`gateway/internal/kafka/consumer/consumer.go:Start`):
     - Подписывается на топик `message.delivered`
     - Читает уведомление
     - Ищет socketID получателя в Redis: `socket:<user_id>`
     - Находит WebSocket соединение
     - Отправляет сообщение через WebSocket

#### PUT /message/upd — Заглушка

**Ответ:** 200 OK

#### GET /message/:id — Заглушка

**Ответ:** 200 OK

#### GET /message/chat/:chat_id — Получить сообщения чата

**Процесс выполнения:**

1. **Handler** (`gateway/internal/handler/message.go:GetChatMessages`):
   - Извлекает `chat_id` из параметров URL
   - Вызывает `h.service.Message.GetMessagesByChatID(chatID)`

2. **Service** (gateway):
   - Вызывает message service

3. **Message Service Repository**:
   - Выполняет SQL запрос:
     ```sql
     SELECT id, chat_id, sender_id, content, sent_at, is_delivered 
     FROM messages 
     WHERE chat_id = $1 
     ORDER BY sent_at ASC
     ```

**Ответ (Response):**
```json
[
    {
        "id": "uuid",
        "chat_id": "uuid",
        "sender_id": "uuid",
        "content": "текст",
        "sent_at": 1234567890,
        "is_delivered": true
    },
    ...
]
```

---

## Frontend (Порт 3000)

Веб-интерфейс на React + TypeScript.

### Компоненты

- **App.tsx** — главный компонент с роутингом
- **Login.tsx** — страница входа
- **Register.tsx** — страница регистрации
- **ChatList.tsx** — список чатов
- **Chat.tsx** — окно чата с сообщениями

### Nginx Конфигурация

```nginx
location /api/ {
    proxy_pass http://gateway:8080/;
}

location /ws {
    proxy_pass http://gateway:8080/ws;
    proxy_http_version 1.1;
    proxy_set_header Upgrade $http_upgrade;
    proxy_set_header Connection "upgrade";
}
```

### Поток данных (Frontend)

1. Пользователь вводит логин/пароль
2. Frontend отправляет POST `/api/auth/login`
3. Сохраняет token в localStorage
4. При каждом запросе добавляет заголовок `Authorization: Bearer <token>`
5. Для WebSocket: подключается к `ws://host/ws?token=<token>`

---

## Схема взаимодействия сервисов

### Регистрация

```
Frontend -> Gateway (/auth/register) -> Auth Service -> PostgreSQL (users)
                                                                  |
Frontend <- Gateway <- Auth Service <-----------------------------+
```

### Вход

```
Frontend -> Gateway (/auth/login) -> Auth Service -> PostgreSQL (users)
                                                                  |
                                                                  v
Frontend <- Gateway <- Auth Service <-- JWT (RS256) <------------+
```

### Создание чата

```
Frontend -> Gateway (/chat/private) -> Chat Service -> PostgreSQL (chats, chat_participants)
                                                                   |
Frontend <- Gateway <- Chat Service <-----------------------------+
```

### Отправка сообщения (WebSocket)

```
User A (WebSocket) -> Gateway -> Kafka (message.send) -> Message Service
                                                             |
                                                             v
                                                       PostgreSQL (messages)
                                                             |
                                                             v
                                    Kafka (message.delivered) <---------+
                                     |                                  |
                                     v                                  |
                               Gateway (Consumer)                      |
                                     |                                  |
                                     v                                  |
                          User B (WebSocket) <-------------------------+
```

---

## Таблицы базы данных

### users

```sql
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
```

### chats

```sql
CREATE TABLE chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    is_private BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
```

### chat_participants

```sql
CREATE TABLE chat_participants (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);
```

### messages

```sql
CREATE TABLE messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_delivered BOOLEAN DEFAULT false,
    sent_at TIMESTAMPTZ DEFAULT now()
);
```

---

## Кэширование (Redis)

### Ключи

| Ключ | Значение | TTL | Описание |
|------|----------|-----|----------|
| `socket:<user_id>` | socketID | 30 мин | Маппинг user -> socket |
| `socket_to_user:<socketID>` | userID | 30 мин | Маппинг socket -> user |

### Использование

1. При подключении WebSocket создаётся маппинг
2. При отправке сообщения проверяется активность соединения
3. При получении сообщения из Kafka ищется socket получателя
