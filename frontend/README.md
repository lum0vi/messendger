# Messenger Frontend

Фронтенд для Messenger приложения на React + TypeScript.

## Технологии

- React 18
- TypeScript
- React Router v6
- Tailwind CSS
- Vite
- WebSocket для real-time сообщений

## Разработка

### Установка зависимостей

```bash
npm install
```

### Запуск в режиме разработки

```bash
npm run dev
```

Приложение будет доступно на http://localhost:3000

### Сборка для production

```bash
npm run build
```

### Предпросмотр production сборки

```bash
npm run preview
```

## Функционал

### Аутентификация
- Регистрация нового пользователя
- Вход в систему
- Автоматический выход при истечении токена

### Пользователи
- Просмотр профиля
- Редактирование профиля (имя, email, пароль)
- Список всех пользователей
- Поиск пользователя по ID
- Поиск пользователя по username

### Чаты
- Список чатов
- Создание приватного чата (1-на-1)
- Создание публичного чата (групповой)
- Просмотр участников чата

### Сообщения
- Отправка сообщений через WebSocket
- Real-time получение новых сообщений
- История сообщений в чате

## Структура проекта

```
src/
├── api/           # API клиенты
├── components/    # React компоненты
│   └── ui/        # UI компоненты (Button, Input, Card)
├── context/       # React контексты
├── hooks/         # Кастомные хуки
├── pages/         # Страницы приложения
└── types/         # TypeScript типы
```

## Docker

### Сборка образа

```bash
docker build -t messenger-frontend .
```

### Запуск контейнера

```bash
docker run -p 3000:80 messenger-frontend
```

## Переменные окружения

Приложение использует проксирование через Vite для API запросов:
- `/api/*` -> `http://localhost:8080/*`
- `/ws` -> `ws://localhost:8080/ws`
