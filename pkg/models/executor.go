package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ExecutorConfig represents the configuration for a task executor
type ExecutorConfig struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	Name         string             `bson:"name"`
	Enabled      bool               `bson:"enabled"`
	WriteConcern WriteConcern       `bson:"write_concern"`
	RetryPolicy  RetryPolicy        `bson:"retry_policy"`
	DLQConfig    DLQConfig          `bson:"dlq_config"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
}

// WriteConcern represents MongoDB write concern settings
type WriteConcern struct {
	Level WriteConcernLevel `bson:"level"`
}

type WriteConcernLevel string

const (
	WriteConcernReplicaAcknowledged WriteConcernLevel = "replica_acknowledged"
	WriteConcernMajority            WriteConcernLevel = "majority"
	WriteConcernUnacknowledged      WriteConcernLevel = "unacknowledged"
	WriteConcernJournaled           WriteConcernLevel = "journaled"
)

// RetryPolicy defines how failed tasks should be retried
type RetryPolicy struct {
	Type        RetryPolicyType `bson:"type"`
	MaxAttempts int             `bson:"max_attempts"` // 0 means unlimited
	Interval    time.Duration   `bson:"interval"`
}

type RetryPolicyType string

const (
	RetryPolicyConstant    RetryPolicyType = "constant"
	RetryPolicyLinear      RetryPolicyType = "linear"
	RetryPolicyExponential RetryPolicyType = "exponential"
)

// DLQConfig represents Dead Letter Queue configuration
type DLQConfig struct {
	Enabled   bool   `bson:"enabled"`
	QueueName string `bson:"queue_name"`
}

// Task represents a single task in the system
type Task struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ExecutorName string             `bson:"executor_name"`
	Status       TaskStatus         `bson:"status"`
	Data         []byte             `bson:"data"`
	Metadata     map[string]string  `bson:"metadata"`
	Error        string             `bson:"error,omitempty"`
	RetryCount   int                `bson:"retry_count"`
	CreatedAt    time.Time          `bson:"created_at"`
	UpdatedAt    time.Time          `bson:"updated_at"`
	StartedAt    *time.Time         `bson:"started_at,omitempty"`
	CompletedAt  *time.Time         `bson:"completed_at,omitempty"`
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"
	TaskStatusInProgress TaskStatus = "in_progress"
	TaskStatusCompleted  TaskStatus = "completed"
	TaskStatusFailed     TaskStatus = "failed"
	TaskStatusDLQ        TaskStatus = "dlq"
)

// TaskProcessor defines the interface that all task processors must implement
type TaskProcessor interface {
	// ProcessTask processes a single task
	ProcessTask(task *Task) error
	// GetTaskSchema returns the JSON schema for the task data
	GetTaskSchema() string
}

// TaskProcessorFactory creates new task processors
type TaskProcessorFactory interface {
	// CreateProcessor creates a new processor for the given executor
	CreateProcessor(config *ExecutorConfig) (TaskProcessor, error)
}
