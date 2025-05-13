# syntax=docker/dockerfile:1
FROM golang:1.22.12-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o manager ./cmd/manager

FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/manager ./manager
COPY frontend ./frontend
EXPOSE 8080
CMD ["./manager"] 