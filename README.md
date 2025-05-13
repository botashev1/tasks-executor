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
- Docker и docker-compose

### Запуск через Docker

Для запуска проекта (backend + MongoDB + frontend):

```bash
git clone https://github.com/yourusername/tasks-executor.git
cd tasks-executor
docker compose up --build -d
```

- Backend будет доступен на: [http://localhost:8080/](http://localhost:8080/)
- MongoDB: `localhost:27017`

### Веб-интерфейс

- Откройте в браузере: [http://localhost:8080/](http://localhost:8080/)
- Или напрямую: [http://localhost:8080/admin_page.html](http://localhost:8080/admin_page.html)

Веб-интерфейс работает через REST API по префиксу `/api/v1`.

## Использование SDK (Go)

SDK позволяет реализовать собственные обработчики задач. Пример:

```go
package main

import (
    "github.com/botashev/tasks-executor/pkg/sdk"
    "github.com/botashev/tasks-executor/pkg/models"
)

type MyTaskHandler struct{}

func (h *MyTaskHandler) ProcessTask(task *models.Task) error {
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

## API

- gRPC и REST API описаны в файле `proto/task_executor.proto`.
- Основные сущности: Executor, Task, DLQ.
- REST API доступен по префиксу `/api/v1`.

## Очистка и массовое наполнение базы

Для очистки коллекции обработчиков используйте:

```bash
docker compose exec mongodb mongosh --eval "db = db.getSiblingDB('task_executor'); db.executors.deleteMany({})"
```

Для массового наполнения базы обработчиками используйте CLI с однострочными JSON-конфигами (см. выше).

## Контейнеризация

- Все сервисы и MongoDB запускаются одной командой через Docker Compose.
- Для разработки удобно использовать volume-монтирование фронтенда: изменения в `frontend/` сразу видны в контейнере.

---

**Tasks Executor** — современная платформа для асинхронной обработки задач с удобным UI, CLI и SDK.

