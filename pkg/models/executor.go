package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

/*
ExecutorConfig represents the configuration for a task executor in the system.
It defines how tasks should be processed, including retry policies, write concerns,
and dead letter queue settings. This configuration is stored in the database and
can be modified at runtime.
*/
type ExecutorConfig struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"` // Unique identifier in the database
	Name         string             `bson:"name"`          // Unique name of the executor
	Enabled      bool               `bson:"enabled"`       // Whether the executor is active
	WriteConcern WriteConcern       `bson:"write_concern"` // Data durability settings
	RetryPolicy  RetryPolicy        `bson:"retry_policy"`  // Task retry configuration
	DLQConfig    DLQConfig          `bson:"dlq_config"`    // Dead letter queue settings
	CreatedAt    time.Time          `bson:"created_at"`    // Creation timestamp
	UpdatedAt    time.Time          `bson:"updated_at"`    // Last update timestamp
}

/*
WriteConcern defines the durability requirements for task operations.
It specifies how many replicas must acknowledge a write operation before it is considered successful.
*/
type WriteConcern struct {
	Level WriteConcernLevel `bson:"level"` // The required acknowledgment level
}

type WriteConcernLevel string

const (
	WriteConcernReplicaAcknowledged WriteConcernLevel = "replica_acknowledged" // At least one replica must acknowledge
	WriteConcernMajority            WriteConcernLevel = "majority"             // Majority of replicas must acknowledge
	WriteConcernUnacknowledged      WriteConcernLevel = "unacknowledged"       // No acknowledgment required
	WriteConcernJournaled           WriteConcernLevel = "journaled"            // Write must be journaled
)

/*
RetryPolicy defines how failed tasks should be retried.
It specifies the retry strategy, maximum number of attempts, and delay between retries.
*/
type RetryPolicy struct {
	Type        RetryPolicyType `bson:"type"`         // The retry strategy to use
	MaxAttempts int             `bson:"max_attempts"` // Maximum number of retry attempts
	Interval    time.Duration   `bson:"interval"`     // Base delay between retries
}

type RetryPolicyType string

const (
	RetryPolicyConstant    RetryPolicyType = "constant"    // Fixed delay between retries
	RetryPolicyLinear      RetryPolicyType = "linear"      // Linearly increasing delay
	RetryPolicyExponential RetryPolicyType = "exponential" // Exponentially increasing delay
)

/*
DLQConfig defines the configuration for the Dead Letter Queue.
The Dead Letter Queue is used to store tasks that have failed after all retry attempts.
*/
type DLQConfig struct {
	Enabled   bool   `bson:"enabled"`    // Whether DLQ is enabled for this executor
	QueueName string `bson:"queue_name"` // Name of the DLQ collection
}

/*
Task represents a unit of work to be processed by an executor.
It contains the task data, metadata, and state information.
*/
type Task struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`          // Unique identifier in the database
	ExecutorName string             `bson:"executor_name"`          // Name of the executor that should process this task
	Status       TaskStatus         `bson:"status"`                 // Current state of the task
	Data         []byte             `bson:"data"`                   // Task payload (JSON)
	Metadata     map[string]string  `bson:"metadata"`               // Additional task metadata
	Error        string             `bson:"error,omitempty"`        // Error message if task failed
	RetryCount   int                `bson:"retry_count"`            // Number of retry attempts
	CreatedAt    time.Time          `bson:"created_at"`             // Creation timestamp
	UpdatedAt    time.Time          `bson:"updated_at"`             // Last update timestamp
	StartedAt    *time.Time         `bson:"started_at,omitempty"`   // When processing started
	CompletedAt  *time.Time         `bson:"completed_at,omitempty"` // When processing completed
}

type TaskStatus string

const (
	TaskStatusPending    TaskStatus = "pending"     // Task is waiting to be processed
	TaskStatusInProgress TaskStatus = "in_progress" // Task is currently being processed
	TaskStatusCompleted  TaskStatus = "completed"   // Task was successfully processed
	TaskStatusFailed     TaskStatus = "failed"      // Task processing failed
	TaskStatusDLQ        TaskStatus = "dlq"         // Task was moved to Dead Letter Queue
)

/*
TaskProcessor defines the interface for processing individual tasks.
Implementations of this interface are responsible for the actual business logic
of task processing and providing a schema for task data validation.
*/
type TaskProcessor interface {
	ProcessTask(task *Task) error
	GetTaskSchema() string
}

/*
TaskProcessorFactory is responsible for creating task processor instances.
It allows for dynamic creation of processors based on executor configuration,
enabling different processing strategies for different executors.
*/
type TaskProcessorFactory interface {
	CreateProcessor(config *ExecutorConfig) (TaskProcessor, error)
}
