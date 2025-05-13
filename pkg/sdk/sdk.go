package sdk

import (
	"github.com/yourusername/tasks-executor/pkg/models"
)

// TaskProcessor is an interface for custom task handlers
// Implement this interface for your task handler
// Register the handler using RegisterProcessor

type TaskProcessor interface {
	ProcessTask(task *models.Task) error
	GetTaskSchema() string
}

var registry = map[string]TaskProcessor{}

// RegisterProcessor registers a task handler by name
func RegisterProcessor(name string, processor TaskProcessor) {
	registry[name] = processor
}

// GetProcessor retrieves a task handler by name
func GetProcessor(name string) (TaskProcessor, bool) {
	p, ok := registry[name]
	return p, ok
}

// Пример пользовательского обработчика
// type MyTaskProcessor struct{}
// func (p *MyTaskProcessor) ProcessTask(task *models.Task) error { /* ... */ return nil }
// func (p *MyTaskProcessor) GetTaskSchema() string { return `{"type":"object"}` }
//
// func init() {
//     RegisterProcessor("my_processor", &MyTaskProcessor{})
// }
