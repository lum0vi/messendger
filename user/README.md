# User Service

## Описание

User Service — микросервис управления пользователями. Предоставляет API для получения и обновления данных пользователей.

**Порт:** 8082

## Требования

Все endpoints требуют заголовок `id` с UUID пользователя (из JWT токена, через Gateway).

## HTTP Endpoints

### GET /me — Получить данные текущего пользователя

**Заголовки:**
```
id: uuid-пользователя
Authorization: Bearer <token>
```

**Процесс обработки:**

1. **Handler** (`handler/user.go:GetMe`)
   - Извлекает `id` из заголовка запроса
   - Вызывает `svc.User.GetMe(id)`

2. **Service**
   - Вызывает репозиторий для получения пользователя

3. **Repository**
   - Выполняет SQL:
     ```sql
     SELECT id, username, email, password_hash, created_at, updated_at 
     FROM users WHERE id = $1
     ```

**Ответ:**
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

### PUT /me — Обновить профиль пользователя

**Заголовки:**
```
id: uuid-пользователя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные:**
```json
{
    "username": "string",  // Опционально
    "email": "string",     // Опционально
    "password": "string"   // Опционально, будет хэшировано
}
```

**Процесс обработки:**

1. **Handler** (`handler/user.go:UpdateMe`)
   - Извлекает `id` из заголовка
   - Десериализует JSON в `models.UpdateMeRequest`
   - Вызывает `svc.User.UpdateMe(id, &req)`

2. **Service**
   - Если передан новый пароль — хэширует его через bcrypt
   - Вызывает репозиторий для обновления

3. **Repository**
   - Выполняет SQL (поля обновляются только если переданы):
     ```sql
     UPDATE users SET 
       username = COALESCE($2, username),
       email = COALESCE($3, email),
       password_hash = COALESCE($4, password_hash),
       updated_at = NOW()
     WHERE id = $1
     ```

**Ответ:**
```json
{}
```

---

### GET /users — Получить всех пользователей

**Заголовки:**
```
id: uuid-пользователя
Authorization: Bearer <token>
```

**Процесс обработки:**

1. **Handler** (`handler/user.go:GetUsers`)
   - Вызывает `svc.User.GetUsers()`

2. **Repository**
   - Выполняет SQL:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users
     ```

**Ответ:**
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

### POST /id — Получить пользователя по ID

**Заголовки:**
```
id: uuid-запрашивающего пользователя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные:**
```json
{
    "id": "uuid-пользователя"
}
```

**Процесс обработки:**

1. **Handler** (`handler/user.go:GetUserById`)
   - Десериализует JSON в `models.GetUserByIDRequest`
   - Вызывает `svc.User.GetUserByID(&req)`

2. **Repository**
   - Выполняет SQL:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users WHERE id = $1
     ```

**Ответ:**
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

### POST /name — Получить пользователя по имени

**Заголовки:**
```
id: uuid-запрашивающего пользователя
Content-Type: application/json
Authorization: Bearer <token>
```

**Входящие данные:**
```json
{
    "username": "string"
}
```

**Процесс обработки:**

1. **Handler** (`handler/user.go:GetUserByName`)
   - Десериализует JSON в `models.GetUserByUsernameRequest`
   - Вызывает `svc.User.GetUserByUsername(&req)`

2. **Repository**
   - Выполняет SQL:
     ```sql
     SELECT id, username, email, created_at, updated_at 
     FROM users WHERE username = $1
     ```

**Ответ:**
```json
{
    "id": "uuid",
    "username": "string",
    "email": "string",
    "created_at": "timestamp",
    "updated_at": "timestamp"
}
```

## Структура проекта

```
user/
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

- **PostgreSQL** — чтение/запись данных пользователей

## Взаимодействие с другими сервисами

```
User Service
    │
    └──→ PostgreSQL (users table)
             └── Получение/обновление данных пользователей
```

## Особенности

- Не хранит пароли — только возвращает их хэши при GET /me
- Обновление полей происходит через COALESCE — переданы только указанные поля
- Требует заголовок `id` от Gateway для идентификации текущего пользователя

## Логи

```bash
docker compose logs -f user-service
```
