package storage

import (
	"context"

	"github.com/botashev/tasks-executor/pkg/models"
)

type Storage interface {
	// Executor operations
	CreateExecutor(ctx context.Context, config *models.ExecutorConfig) error
	UpdateExecutor(ctx context.Context, config *models.ExecutorConfig) error
	GetExecutor(ctx context.Context, name string) (*models.ExecutorConfig, error)
	ListExecutors(ctx context.Context) ([]*models.ExecutorConfig, error)
	DeleteExecutor(ctx context.Context, name string) error

	// Task operations
	AddTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, id string) (*models.Task, error)
	UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus, error string) error
	GetNextTask(ctx context.Context, executorName string) (*models.Task, error)
	MoveToDLQ(ctx context.Context, task *models.Task) error
	GetDLQTasks(ctx context.Context, executorName string) ([]*models.Task, error)
	ClearDLQ(ctx context.Context, executorName string) error
}

type StorageConfig struct {
	MongoURI      string
	Database      string
	ExecutorsColl string
	TasksColl     string
	DLQColl       string
}
