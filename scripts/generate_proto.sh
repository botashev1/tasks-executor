#!/bin/bash

# Add Go bin to PATH
export PATH="$PATH:$(go env GOPATH)/bin"

# Install plugins if not present
if ! command -v protoc-gen-go &> /dev/null; then
  go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi
if ! command -v protoc-gen-go-grpc &> /dev/null; then
  go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi
if ! command -v protoc-gen-grpc-gateway &> /dev/null; then
  go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
fi

# Создаем директории для сгенерированных файлов
mkdir -p proto/taskexecutor

# Генерируем код из proto-файлов
protoc -I. -Ithird_party --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative --grpc-gateway_out=. --grpc-gateway_opt=paths=source_relative proto/task_executor.proto

# Генерируем код из примера заказа
protoc --go_out=. --go_opt=paths=source_relative \
    examples/order.proto

echo "Proto files generated successfully" 