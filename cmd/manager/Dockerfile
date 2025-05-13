FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем файлы зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o manager ./cmd/manager

# Финальный образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл
COPY --from=builder /app/manager .

# Копируем файлы фронтенда
COPY --from=builder /app/frontend ./frontend

# Запускаем приложение
CMD ["./manager"] 