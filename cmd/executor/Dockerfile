FROM golang:1.22-alpine AS builder

WORKDIR /app

# Копируем только файлы зависимостей
COPY go.mod go.sum ./

# Скачиваем зависимости с таймаутом
RUN go mod download -x

# Копируем исходный код
COPY . .

# Собираем приложение
RUN cd cmd/executor && go build -o /app/executor

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/executor .
CMD ["./executor"] 