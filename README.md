# Messenger Application

Полнофункциональное приложение для обмена сообщениями с микросервисной архитектурой.

## Содержание

- [Архитектура](#архитектура)
- [Требования](#требования)
- [Быстрый старт](#быстрый-старт)
- [Swagger документация](#swagger-документация)
- [Сервисы](#сервисы)
- [Миграции базы данных](#миграции-базы-данных)
- [Логи](#логи)
- [Переменные окружения](#переменные-окружения)
- [API Endpoints](#api-endpoints)
- [Тестирование](#тестирование)
- [Устранение проблем](#устранение-проблем)

## Архитектура

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│   Frontend  │────▶│   Gateway   │────▶│ Auth Service│
│   (React)   │     │   (Gin)     │     │   (Gin)     │
└─────────────┘     └─────────────┘     └─────────────┘
                           │                    │
                           ▼                    ▼
                    ┌─────────────┐     ┌─────────────┐
                    │    Kafka    │     │  PostgreSQL │
                    │             │     │             │
                    └─────────────┘     └─────────────┘
                           │                    │
                           ▼                    ▼
                    ┌─────────────┐     ┌─────────────┐
                    │   Message   │     │  User/Chat  │
                    │  Service    │     │  Services   │
                    └─────────────┘     └─────────────┘
```

## Требования

- Docker и Docker Compose
- OpenSSL (для генерации ключей)
- 4GB RAM минимум

## Быстрый старт

### Полный запуск приложения

```bash
# 1. Переходим в директорию gateway
cd gateway

# 2. Запускаем все сервисы
docker compose up -d

# 3. Ждём пока PostgreSQL станет здоровым (зеленым)
docker compose ps

# 4. Создаём таблицы в базе данных
docker exec postgres psql -U user -d postgres -f /dev/stdin <<'EOF'
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    is_private BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chat_participants (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_delivered BOOLEAN DEFAULT false,
    sent_at TIMESTAMPTZ DEFAULT now()
);
EOF

# 5. Проверяем доступность
# - Веб-интерфейс: http://localhost:3000
# - API: http://localhost:8080
# - Swagger: http://localhost:8080/swagger/index.html
```

### Пошаговая инструкция

#### 1. Клонирование и запуск

```bash
cd gateway
docker compose up -d
```

### 2. Проверка статуса сервисов

```bash
cd gateway
docker compose ps
```

Дождитесь, пока PostgreSQL и Redis покажут статус `healthy` (зеленым).

### 3. Создание таблиц в базе данных

```bash
docker exec postgres psql -U user -d postgres -f /dev/stdin <<'EOF'
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chats (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255),
    is_private BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

CREATE TABLE IF NOT EXISTS chat_participants (
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    joined_at TIMESTAMPTZ DEFAULT now(),
    PRIMARY KEY (chat_id, user_id)
);

CREATE TABLE IF NOT EXISTS messages (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id UUID REFERENCES chats(id) ON DELETE CASCADE,
    sender_id UUID REFERENCES users(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    is_delivered BOOLEAN DEFAULT false,
    sent_at TIMESTAMPTZ DEFAULT now()
);
EOF
```

### 4. Доступ к приложению

- **Веб-интерфейс**: http://localhost:3000
- **API Gateway**: http://localhost:8080
- **Swagger**: http://localhost:8080/swagger/index.html

### 5. Swagger документация

После запуска приложения откройте в браузере:

```
http://localhost:8080/swagger/index.html
```

Здесь вы найдёте полную документацию всех API эндпоинтов с возможностью тестирования прямо из браузера.

**Особенности Swagger в этом проекте:**
- Автоматическая генерация из аннотаций Go кода
- Тестирование запросов прямо в браузере
- Примеры запросов и ответов

**Для защищённых эндпоинтов** (требующих JWT):
1. Сначала выполните Login через `/auth/login`
2. Скопируйте полученный токен
3. В Swagger нажмите кнопку "Authorize"
4. Введите токен в формате: `Bearer <ваш_токен>`

## Сервисы

| Сервис | Порт | Описание |
|--------|------|----------|
| Gateway | 8080 | API Gateway |
| Auth Service | 8081 | Аутентификация |
| User Service | 8082 | Пользователи |
| Chat Service | 8083 | Чаты |
| Message Service | 8084 | Сообщения |
| PostgreSQL | 5432 | База данных |
| Redis | 6379 | Кэш |
| Kafka | 9092 | Сообщения |
| Frontend | 3000 | Веб-интерфейс |

## Логи

```bash
# Все сервисы
docker compose logs -f

# Конкретный сервис
docker compose logs -f auth-service
docker compose logs -f gateway
```

## API Endpoints

### Аутентификация
- POST /auth/register - Регистрация
- POST /auth/login - Вход

### Пользователи (требуется JWT)
- GET /user/me - Текущий пользователь
- PUT /user/me - Обновить профиль
- GET /user/users - Все пользователи
- POST /user/id - Пользователь по ID
- POST /user/name - Пользователь по имени

### Чаты (требуется JWT)
- GET /chat/ - Мои чаты
- POST /chat/private - Приватный чат
- POST /chat/public - Публичный чат
- GET /chat/:chat_id/users - Участники

### Сообщения (требуется JWT)
- PUT /message/upd - Статус сообщения
- GET /message/:id - Недоставленные

### WebSocket
- GET /ws - Подключение

## Тестирование

Импортируйте `messenger.postman_collection.json` в Postman.

```bash
# Регистрация
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass","email":"a@b.c"}'

# Логин
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"user","password":"pass"}'
```

## Устранение проблем

```bash
# Перезапуск
docker compose down && docker compose up -d

# Очистка томов
docker compose down -v && docker compose up -d

# Пересборка
docker compose build --no-cache
```
