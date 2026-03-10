# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, чтобы сначала установить зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь исходный код
COPY . .

# Собираем бинарный файл
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/main.go

# Этап минимального запуска
FROM alpine:3.19

# Устанавливаем необходимые зависимости
RUN apk add --no-cache ca-certificates

# Создаем непривилегированного пользователя
RUN adduser -D -g '' appuser

# Копируем бинарник из предыдущего этапа
COPY --from=builder /app/main /main

# Меняем пользователя
USER appuser
EXPOSE 8082

# Устанавливаем команду по умолчанию
ENTRYPOINT ["/main"]