package executors

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/botashev/tasks-executor/pkg/manager"
	pb "github.com/botashev/tasks-executor/proto"
)

// ExampleProcessor представляет пример обработчика задач
type ExampleProcessor struct {
	manager *manager.Manager
}

// NewExampleProcessor создает новый экземпляр ExampleProcessor
func NewExampleProcessor(m *manager.Manager) *ExampleProcessor {
	return &ExampleProcessor{
		manager: m,
	}
}

// Register регистрирует обработчик в системе
func (p *ExampleProcessor) Register() error {
	return p.manager.RegisterExecutor("example_processor")
}

// ProcessTask обрабатывает входящую задачу
func (p *ExampleProcessor) ProcessTask(task *pb.Task) error {
	log.Printf("Processing task %s", task.Id)

	// Пример обработки данных задачи
	var data map[string]interface{}
	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal task data: %v", err)
	}

	// Пример логики обработки
	log.Printf("Task data: %+v", data)

	// Здесь можно добавить любую бизнес-логику
	// Например, отправку уведомлений, обработку данных и т.д.

	// Имитация обработки
	time.Sleep(1 * time.Second)

	return nil
}

// GetTaskSchema возвращает JSON-схему для валидации данных задачи
func (p *ExampleProcessor) GetTaskSchema() string {
	return `{
		"type": "object",
		"properties": {
			"message": {
				"type": "string",
				"description": "Сообщение для обработки"
			},
			"priority": {
				"type": "integer",
				"description": "Приоритет задачи",
				"minimum": 1,
				"maximum": 10
			}
		},
		"required": ["message"]
	}`
}

// Start начинает обработку задач
func (p *ExampleProcessor) Start() error {
	for {
		// Получаем следующую задачу
		task, err := p.manager.GetNextTask("example_processor")
		if err != nil {
			log.Printf("Error getting next task: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if task == nil {
			// Нет задач для обработки
			time.Sleep(1 * time.Second)
			continue
		}

		// Обрабатываем задачу
		err = p.ProcessTask(task)
		if err != nil {
			log.Printf("Error processing task %s: %v", task.Id, err)
			// Обновляем статус задачи на FAILED
			if err := p.manager.UpdateTaskStatus(task.Id, pb.TaskStatus_TASK_STATUS_FAILED, err.Error()); err != nil {
				log.Printf("Error updating task status: %v", err)
			}
			continue
		}

		// Обновляем статус задачи на COMPLETED
		if err := p.manager.UpdateTaskStatus(task.Id, pb.TaskStatus_TASK_STATUS_COMPLETED, ""); err != nil {
			log.Printf("Error updating task status: %v", err)
		}
	}
}
