package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type RetryPolicy struct {
	Type        RetryPolicyType `bson:"type"`
	MaxAttempts int             `bson:"max_attempts"`
	Interval    time.Duration   `bson:"interval"`
}

type RetryPolicyType string

const (
	RetryPolicyConstant    RetryPolicyType = "constant"
	RetryPolicyLinear      RetryPolicyType = "linear"
	RetryPolicyExponential RetryPolicyType = "exponential"
)

type DLQConfig struct {
	Enabled   bool   `bson:"enabled"`
	QueueName string `bson:"queue_name"`
}

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

type TaskProcessor interface {
	ProcessTask(task *Task) error
	GetTaskSchema() string
}

type TaskProcessorFactory interface {
	CreateProcessor(config *ExecutorConfig) (TaskProcessor, error)
}
