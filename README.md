# Tasks Executor

Гибкая система асинхронной обработки задач с возможностью регистрации пользовательских обработчиков через веб-интерфейс и SDK.

## Назначение

Tasks Executor предоставляет инфраструктуру для надёжного и масштабируемого выполнения задач в распределённых системах. Система обеспечивает автоматические повторные попытки, горизонтальное масштабирование обработчиков, мониторинг и изоляцию.

## Архитектура

- **Backend** (Go):
  - gRPC и REST API (через gRPC-Gateway)
  - MongoDB для хранения задач и обработчиков
  - Отдача статического frontend
- **Frontend** (HTML+JS):
  - Веб-интерфейс администратора (управление обработчиками, мониторинг)
  - Реализован как статика (без сборщиков)
- **SDK** (Go):
  - Для реализации пользовательских обработчиков задач

### Структура репозитория

- `cmd/manager/` — основной backend-сервис
- `frontend/` — статический веб-интерфейс
  - `admin_page.html` — главная страница
  - `assets/js/` — логика работы интерфейса и API
  - `assets/styles/` — стили
  - `components/` — HTML-компоненты интерфейса
- `pkg/` — бизнес-логика, API, SDK, модели, хранилище
- `proto/` — gRPC/REST API спецификации
- `Dockerfile`, `docker-compose.yml` — контейнеризация и быстрый запуск

## Быстрый старт

### Требования
- Go 1.22+
- Docker и docker-compose (для быстрого старта)

### Запуск через Docker

Для быстрого запуска всего проекта (backend + MongoDB) используйте Docker и docker-compose:

```bash
git clone https://github.com/yourusername/tasks-executor.git
cd tasks-executor
docker-compose up --build -d
```

- Backend будет доступен на: [http://localhost:8080/](http://localhost:8080/)
- MongoDB: `localhost:27017`

### Локальный запуск backend

Если хотите запускать backend локально (без Docker):

```bash
go mod download
go build -o manager ./cmd/manager
./manager
```

> **MongoDB** должна быть запущена отдельно (например, через Docker или как сервис).

### Веб-интерфейс

- Откройте в браузере: [http://localhost:8080/](http://localhost:8080/)
- Или напрямую: [http://localhost:8080/admin_page.html](http://localhost:8080/admin_page.html)

Веб-интерфейс работает через REST API по префиксу `/api/v1`.

## Использование SDK (Go)

SDK позволяет реализовать собственные обработчики задач. Пример:

```go
package main

import (
    "context"
    "github.com/yourusername/tasks-executor/pkg/sdk"
)

type MyTaskHandler struct{}

func (h *MyTaskHandler) ProcessTask(task *sdk.Task) error {
    // Ваша логика обработки задачи
    return nil
}

func (h *MyTaskHandler) GetTaskSchema() string {
    return `{"type":"object"}`
}

func main() {
    sdk.RegisterProcessor("my_handler", &MyTaskHandler{})
    // ... запуск сервиса
}
```

