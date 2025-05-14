package storage

import (
	"context"

	"github.com/botashev/tasks-executor/pkg/models"
)

/*
Storage defines the interface for persistent storage operations in the task execution system.
This interface provides methods for managing both executors and tasks, including their lifecycle
and state transitions. Implementations of this interface should ensure thread safety and
proper handling of concurrent operations.

The interface is divided into two main sections:
1. Executor operations - for managing executor configurations
2. Task operations - for managing task lifecycle and state
*/
type Storage interface {
	// Executor operations
	/*
		CreateExecutor adds a new executor configuration to the storage.
		Returns an error if an executor with the same name already exists.
	*/
	CreateExecutor(ctx context.Context, config *models.ExecutorConfig) error

	/*
		UpdateExecutor modifies an existing executor configuration.
		Returns an error if the executor doesn't exist.
	*/
	UpdateExecutor(ctx context.Context, config *models.ExecutorConfig) error

	/*
		GetExecutor retrieves an executor configuration by its name.
		Returns nil and an error if the executor doesn't exist.
	*/
	GetExecutor(ctx context.Context, name string) (*models.ExecutorConfig, error)

	/*
		ListExecutors returns all registered executor configurations.
		Returns an empty slice if no executors are found.
	*/
	ListExecutors(ctx context.Context) ([]*models.ExecutorConfig, error)

	/*
		DeleteExecutor removes an executor configuration by its name.
		Returns an error if the executor doesn't exist.
	*/
	DeleteExecutor(ctx context.Context, name string) error

	// Task operations
	/*
		AddTask creates a new task in the storage.
		The task will be in PENDING state initially.
	*/
	AddTask(ctx context.Context, task *models.Task) error

	/*
		GetTask retrieves a task by its ID.
		Returns nil and an error if the task doesn't exist.
	*/
	GetTask(ctx context.Context, id string) (*models.Task, error)

	/*
		UpdateTaskStatus changes the status of a task and optionally sets an error message.
		This method should also update the task's timestamps based on the new status.
	*/
	UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus, error string) error

	/*
		GetNextTask retrieves the next available task for an executor.
		The task should be in PENDING state and not assigned to any other executor.
		Returns nil if no tasks are available.
	*/
	GetNextTask(ctx context.Context, executorName string) (*models.Task, error)

	/*
		MoveToDLQ moves a failed task to the Dead Letter Queue.
		This is typically called when a task has exceeded its retry attempts.
	*/
	MoveToDLQ(ctx context.Context, task *models.Task) error

	/*
		GetDLQTasks retrieves all tasks in the Dead Letter Queue for an executor.
		Returns an empty slice if no tasks are in the DLQ.
	*/
	GetDLQTasks(ctx context.Context, executorName string) ([]*models.Task, error)

	/*
		ClearDLQ removes all tasks from the Dead Letter Queue for an executor.
		This operation cannot be undone.
	*/
	ClearDLQ(ctx context.Context, executorName string) error
}

/*
StorageConfig holds the configuration parameters for storage implementations.
This structure is used to initialize storage backends with the necessary
connection details and collection names.
*/
type StorageConfig struct {
	MongoURI      string // MongoDB connection URI
	Database      string // Database name
	ExecutorsColl string // Collection name for executor configurations
	TasksColl     string // Collection name for tasks
	DLQColl       string // Collection name for dead letter queue
}
