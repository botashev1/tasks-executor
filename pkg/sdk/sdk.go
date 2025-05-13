package sdk

import (
	"github.com/botashev/tasks-executor/pkg/models"
)

type TaskProcessor interface {
	ProcessTask(task *models.Task) error
	GetTaskSchema() string
}

var registry = map[string]TaskProcessor{}

func RegisterProcessor(name string, processor TaskProcessor) {
	registry[name] = processor
}

func GetProcessor(name string) (TaskProcessor, bool) {
	p, ok := registry[name]
	return p, ok
}
