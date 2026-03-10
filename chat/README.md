# Chat Service

## Описание

Chat Service — микросервис управления чатами. Создает и управляет приватными и публичными чатами, а также участниками.

**Порт:** 8083

## Требования

Все endpoints требуют заголовок `id` с UUID пользователя (из JWT токена, через Gateway).

## HTTP Endpoints

### POST /chat/private — Создать приватный чат

**Заголовки:**
```
id: uuid-создателя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные:**
```json
{
    "friend_id": "uuid-другого-пользователя"
}
```

**Процесс обработки:**

1. **Handler** (`handler/chat.go:CreatePrivateChat`)
   - Извлекает `id` создателя из заголовка
   - Десериализует JSON в `models.CreatePrivateChatRequest`
   - Вызывает `svc.Chat.CreatePrivateChat(id, req)`

2. **Repository** (`repository/chat.go:CreatePrivateChat`)
   - Начинает транзакцию
   - Создает запись в `chats`:
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

**Ответ:**
```json
{
    "chat_id": "uuid-чата"
}
```

**Коды ответа:**
- 201 — успешно создан
- 400 — некорректные данные
- 401 — не авторизован

---

### POST /chat/public — Создать публичный чат

**Заголовки:**
```
id: uuid-создателя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные:**
```json
{
    "name": "Название чата",
    "participant_id": ["uuid1", "uuid2", ...]
}
```

**Процесс обработки:**

1. **Handler** (`handler/chat.go:CreatePublicChat`)
   - Извлекает `id` создателя из заголовка
   - Десериализует JSON в `models.CreatePublicChatRequest`
   - Вызывает `svc.Chat.CreatePublicChat(id, req)`

2. **Repository** (`repository/chat.go:CreatePublicChat`)
   - Начинает транзакцию
   - Создает запись в `chats`:
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

**Ответ:**
```json
{
    "chat_id": "uuid-чата"
}
```

---

### GET /chat/ — Получить чаты пользователя

**Заголовки:**
```
id: uuid-пользователя
Authorization: Bearer <token>
```

**Процесс обработки:**

1. **Handler** (`handler/chat.go:GetChats`)
   - Извлекает `id` пользователя из заголовка
   - Вызывает `svc.Chat.GetChats(id)`

2. **Repository** (`repository/chat.go:GetChats`)
   - Выполняет SQL:
     ```sql
     SELECT chat_id FROM chat_participants WHERE user_id = $1
     ```

**Ответ:**
```json
{
    "chat_id": ["uuid-чата1", "uuid-чата2", ...]
}
```

---

### GET /chat/:chat_id/users — Получить участников чата

**Параметры URL:**
- `chat_id` — UUID чата

**Заголовки:**
```
id: uuid-пользователя
Authorization: Bearer <token>
```

**Процесс обработки:**

1. **Handler** (`handler/chat.go:GetChatUsers`)
   - Извлекает `chat_id` из параметров URL
   - Вызывает `svc.Chat.GetUsersChat(chatID)`

2. **Repository** (`repository/chat.go:GetUsersChat`)
   - Выполняет SQL:
     ```sql
     SELECT user_id FROM chat_participants WHERE chat_id = $1
     ```

**Ответ:**
```json
{
    "users": ["uuid-user1", "uuid-user2", ...]
}
```

## Структура проекта

```
chat/
├── cmd/main.go              # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── errors/              # Обработка ошибок
│   ├── handler/             # HTTP обработчики
│   ├── middleware/          # Промежуточное ПО
│   ├── models/              # Модели данных
│   ├── repository/          # Работа с БД
│   └── service/             # Бизнес-логика
└── go.mod
```

## Зависимости

- **PostgreSQL** — хранение чатов и участников

## Таблицы БД

### chats
| Поле | Тип | Описание |
|------|-----|----------|
| id | UUID | PK, автогенерируется |
| name | VARCHAR(255) | Название чата (nullable для приватных) |
| is_private | BOOLEAN | true — приватный, false — публичный |
| created_at | TIMESTAMPTZ | Время создания |
| updated_at | TIMESTAMPTZ | Время обновления |

### chat_participants
| Поле | Тип | Описание |
|------|-----|----------|
| chat_id | UUID | FK на chats |
| user_id | UUID | FK на users |
| joined_at | TIMESTAMPTZ | Время присоединения |

## Взаимодействие с другими сервисами

```
Chat Service
    │
    ├──→ PostgreSQL (chats table)
    │        └── Создание/получение чатов
    │
    ├──→ PostgreSQL (chat_participants table)
    │        └── Управление участниками
    │
    └──→ Gateway
             └── Возвращает chat_id для дальнейшей работы
```

## Особенности

- Приватные чаты не имеют названия (name = NULL)
- Создание чата происходит в транзакции — либо все операции успешны, либо ни одной
- При создании приватного чата автоматически добавляются оба пользователя
- При создании публичного чата создатель автоматически добавляется в участники

## Логи

```bash
docker compose logs -f chat-service
```
