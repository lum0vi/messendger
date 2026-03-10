# Auth Service

## Описание

Auth Service — микросервис аутентификации и регистрации пользователей. Отвечает за безопасное хранение учетных данных и генерацию JWT токенов для авторизации.

**Порт:** 8081

## Функции

1. **Регистрация пользователей** — создание новых учетных записей с хэшированием пароля
2. **Аутентификация** — проверка credentials и выдача JWT токена
3. **Генерация JWT** — создание подписанных токенов с использованием RSA-256

## HTTP Endpoints

### POST /auth/register — Регистрация пользователя

**Входящие данные:**
```json
{
    "username": "string",  // Обязательно, уникальное
    "password": "string",  // Обязательно
    "email": "string"      // Обязательно, уникальное
}
```

**Процесс обработки:**

1. **Handler** (`handler/auth.go:Register`)
   - Десериализует JSON в `models.RegisterRequest`
   - Вызывает `srv.Auth.Register(&req)`

2. **Service** (`service/auth.go:Register`)
   - Хэширует пароль с помощью bcrypt: `hash.HashPass(req.Password)`
   - Передает данные в репозиторий

3. **Repository** (`repository/postgres_auth.go:Register`)
   - Выполняет SQL:
     ```sql
     INSERT INTO users (username, password_hash, email) 
     VALUES ($1, $2, $3) RETURNING id
     ```
   - При ошибке уникальности username возвращает ошибку 409

**Ответ:**
```json
{
    "id": "uuid-пользователя"
}
```

**Коды ответа:**
- 200 — успешно
- 400 — некорректные данные
- 409 — username уже занят
- 500 — внутренняя ошибка

---

### POST /auth/login — Вход пользователя

**Входящие данные:**
```json
{
    "username": "string",  // Обязательно
    "password": "string"   // Обязательно
}
```

**Процесс обработки:**

1. **Handler** (`handler/auth.go:Login`)
   - Десериализует JSON в `models.LoginRequest`
   - Вызывает `srv.Auth.Login(&req)`

2. **Service** (`service/auth.go:Login`)
   - Получает id пользователя из репозитория
   - Создает JWT токен: `jwt.CreateJWT(id)`

3. **Repository** (`repository/postgres_auth.go:Login`)
   - Выполняет SQL:
     ```sql
     SELECT id, password_hash FROM users WHERE username = $1
     ```
   - Проверяет пароль: `hash.VerifyPass(req.Password, passHash)`
   - При неверном пароле возвращает ошибку 401

4. **JWT Создание** (`jwt/jwt.go:CreateJWT`)
   - Формирует claims: `{"user_id": "...", "exp": ...}`
   - Подписывает приватным RSA ключом (RS256)

**Ответ:**
```json
{
    "token": "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Коды ответа:**
- 200 — успешно
- 400 — некорректные данные
- 401 — неверный пароль
- 500 — внутренняя ошибка

## Структура проекта

```
auth/
├── cmd/main.go              # Точка входа
├── internal/
│   ├── config/              # Конфигурация
│   ├── errors/              # Обработка ошибок
│   ├── handler/             # HTTP обработчики
│   ├── hash/                # Хэширование паролей (bcrypt)
│   ├── jwt/                 # Создание JWT токенов
│   ├── jwtutil/             # Утилиты для работы с ключами
│   ├── middleware/          # Промежуточное ПО
│   ├── models/              # Модели данных
│   ├── repository/          # Работа с БД
│   └── service/             # Бизнес-логика
└── go.mod
```

## Зависимости

- **PostgreSQL** — хранение пользователей
- **bcrypt** — хэширование паролей
- **RSA ключи** — подпись JWT токенов

## Ключи JWT

Приватный и публичный ключи RSA загружаются из директории `/app/keys/`:
- `private.pem` — для подписи токенов
- `public.pem` — для валидации токенов (используется в Gateway)

## Взаимодействие с другими сервисами

```
Auth Service
    │
    ├──→ PostgreSQL (users table)
    │        └── Сохранение/получение пользователей
    │
    └──→ Gateway
             └── Получает JWT токен для дальнейшей авторизации
```

## Логи

```bash
docker compose logs -f auth-service
```
