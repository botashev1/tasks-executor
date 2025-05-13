package executors

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/botashev/tasks-executor/pkg/manager"
	pb "github.com/botashev/tasks-executor/proto"
)

type ExampleProcessor struct {
	manager *manager.Manager
}

func NewExampleProcessor(m *manager.Manager) *ExampleProcessor {
	return &ExampleProcessor{
		manager: m,
	}
}

func (p *ExampleProcessor) Register() error {
	return p.manager.RegisterExecutor("example_processor")
}

func (p *ExampleProcessor) ProcessTask(task *pb.Task) error {
	log.Printf("Processing task %s", task.Id)

	var data map[string]interface{}
	if err := json.Unmarshal(task.Data, &data); err != nil {
		return fmt.Errorf("failed to unmarshal task data: %v", err)
	}

	log.Printf("Task data: %+v", data)

	time.Sleep(1 * time.Second)

	return nil
}

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

func (p *ExampleProcessor) Start() error {
	for {
		task, err := p.manager.GetNextTask("example_processor")
		if err != nil {
			log.Printf("Error getting next task: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if task == nil {
			time.Sleep(1 * time.Second)
			continue
		}

		err = p.ProcessTask(task)
		if err != nil {
			log.Printf("Error processing task %s: %v", task.Id, err)
			if err := p.manager.UpdateTaskStatus(task.Id, pb.TaskStatus_TASK_STATUS_FAILED, err.Error()); err != nil {
				log.Printf("Error updating task status: %v", err)
			}
			continue
		}

		if err := p.manager.UpdateTaskStatus(task.Id, pb.TaskStatus_TASK_STATUS_COMPLETED, ""); err != nil {
			log.Printf("Error updating task status: %v", err)
		}
	}
}
