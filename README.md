# Tasks Executor

Гибкая система асинхронной обработки задач с возможностью регистрации пользовательских обработчиков через веб-интерфейс и SDK.

## Назначение

Tasks Executor предоставляет инфраструктуру для надёжного и масштабируемого выполнения задач в распределённых системах. Система обеспечивает:

- Автоматические повторные попытки с настраиваемыми стратегиями (постоянная, линейная, экспоненциальная задержка)
- Горизонтальное масштабирование обработчиков
- Dead Letter Queue для обработки неудачных задач
- Мониторинг состояния задач и обработчиков
- Изоляцию обработчиков и их конфигураций
- Настраиваемые уровни гарантий записи (Write Concern)

## Архитектура

Система состоит из трёх основных компонентов:

### Backend (Go)
- gRPC и REST API (через gRPC-Gateway)
- MongoDB для хранения задач, обработчиков и DLQ
- Отдача статического frontend
- Поддержка конфигурации через переменные окружения

### Frontend (HTML+JS)
- Веб-интерфейс администратора
  - Управление обработчиками (создание, редактирование, удаление)
  - Мониторинг состояния задач
  - Просмотр DLQ
- Реализован как статика (без сборщиков)
- Адаптивный дизайн

### SDK (Go)
- Интерфейс для реализации пользовательских обработчиков задач
- Валидация данных задач через JSON Schema
- Регистрация обработчиков в системе

## Запуск

### Требования
- Go 1.22+
- Docker и docker-compose
- MongoDB 6.0+

### Запуск

```bash
git clone https://github.com/yourusername/tasks-executor.git
cd tasks-executor
docker compose up --build -d
```

Сервисы будут доступны:
- Backend API: [http://localhost:8080/](http://localhost:8080/)
- Веб-интерфейс: [http://localhost:8080/admin_page.html](http://localhost:8080/admin_page.html)
- MongoDB: `localhost:27017`
- gRPC: `localhost:50051`

### Конфигурация

```bash
MANAGER_PORT=8080              # Порт для HTTP API
MANAGER_GRPC_PORT=50051        # Порт для gRPC
MONGO_URI=mongodb://localhost:27017  # URI MongoDB
MONGO_DB=task_executor         # Имя базы данных
```

## Использование SDK (Go)

```go
package main

import (
    "encoding/json"
    "github.com/botashev/tasks-executor/pkg/sdk"
    "github.com/botashev/tasks-executor/pkg/models"
)

type MyTaskData struct {
    Message  string `json:"message"`
    Priority int    `json:"priority"`
}

type MyTaskHandler struct{}

func (h *MyTaskHandler) ProcessTask(task *models.Task) error {
    var data MyTaskData
    if err := json.Unmarshal(task.Data, &data); err != nil {
        return err
    }
    // Ваша логика обработки задачи
    return nil
}

func (h *MyTaskHandler) GetTaskSchema() string {
    return `{
        "type": "object",
        "properties": {
            "message": {
                "type": "string",
                "description": "Сообщение для обработки"
            },
            "priority": {
                "type": "integer",
                "minimum": 1,
                "maximum": 10
            }
        },
        "required": ["message"]
    }`
}

func main() {
    sdk.RegisterProcessor("my_handler", &MyTaskHandler{})
    // ... запуск сервиса
}
```

## API

### REST API

Доступно по префиксу `/api/v1`:

- `GET /api/v1/executors` - список обработчиков
- `POST /api/v1/executors` - создание обработчика
- `GET /api/v1/executors/{id}` - информация об обработчике
- `PUT /api/v1/executors/{id}` - обновление обработчика
- `DELETE /api/v1/executors/{id}` - удаление обработчика
- `GET /api/v1/tasks` - список задач
- `POST /api/v1/tasks` - создание задачи
- `GET /api/v1/tasks/{id}` - информация о задаче
- `PUT /api/v1/tasks/{id}/status` - обновление статуса задачи

### gRPC API

Полная спецификация в `proto/task_executor.proto`. Основные сервисы:

- `TaskExecutorManager` - управление обработчиками и задачами
- `TaskExecutor` - выполнение задач

## Управление данными

```bash
# Очистка обработчиков
docker compose exec mongodb mongosh --eval "db = db.getSiblingDB('task_executor'); db.executors.deleteMany({})"

# Очистка задач
docker compose exec mongodb mongosh --eval "db = db.getSiblingDB('task_executor'); db.tasks.deleteMany({})"

# Очистка DLQ
docker compose exec mongodb mongosh --eval "db = db.getSiblingDB('task_executor'); db.dlq.deleteMany({})"
```


