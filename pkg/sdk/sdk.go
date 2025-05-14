package sdk

import (
	"github.com/botashev/tasks-executor/pkg/models"
)

/*
TaskProcessor is the core interface that defines how tasks should be processed in the system.
Implementations of this interface are responsible for:
1. Processing individual tasks according to their specific business logic
2. Providing a JSON schema that defines the expected structure of task data

The interface is designed to be simple yet flexible, allowing for various types of task processors
to be implemented while maintaining a consistent contract with the task execution system.
*/
type TaskProcessor interface {
	// ProcessTask handles the execution of a single task.
	// It should implement the business logic specific to the task type.
	// Returns an error if the task processing fails.
	ProcessTask(task *models.Task) error

	// GetTaskSchema returns a JSON schema string that defines the expected structure
	// of the task's data field. This schema is used for validation before task processing.
	GetTaskSchema() string
}

// registry maintains a mapping of processor names to their implementations
var registry = map[string]TaskProcessor{}

/*
RegisterProcessor adds a new task processor to the system's registry.
This function should be called during application initialization to make
task processors available for use by the task execution system.

Parameters:
- name: A unique identifier for the processor
- processor: An implementation of the TaskProcessor interface
*/
func RegisterProcessor(name string, processor TaskProcessor) {
	registry[name] = processor
}

/*
GetProcessor retrieves a task processor from the registry by its name.
Returns the processor and a boolean indicating whether it was found.

Parameters:
- name: The name of the processor to retrieve

Returns:
- TaskProcessor: The processor implementation if found
- bool: True if the processor was found, false otherwise
*/
func GetProcessor(name string) (TaskProcessor, bool) {
	p, ok := registry[name]
	return p, ok
}
