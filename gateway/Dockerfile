# Этап сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию внутри контейнера
WORKDIR /app

# Копируем go.mod и go.sum, чтобы сначала установить зависимости
COPY go.mod go.sum ./
COPY ./keys ./keys
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
COPY --from=builder /app/keys /app/keys

# Меняем владельца файлов, чтобы appuser имел доступ
RUN chown -R appuser:appuser /app/keys

# Меняем пользователя
USER appuser
EXPOSE 8080

# Устанавливаем команду по умолчанию
ENTRYPOINT ["/main"]