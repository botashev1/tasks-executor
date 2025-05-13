package manager

import (
	"context"
	"fmt"
	"time"

	"github.com/botashev/tasks-executor/pkg/models"
	"github.com/botashev/tasks-executor/pkg/storage"
	pb "github.com/botashev/tasks-executor/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Service struct {
	pb.UnimplementedTaskExecutorManagerServer
	storage storage.Storage
}

func NewService(storage storage.Storage) *Service {
	return &Service{
		storage: storage,
	}
}

// Task Management
func (s *Service) AddTask(ctx context.Context, req *pb.AddTaskRequest) (*pb.AddTaskResponse, error) {
	task := &models.Task{
		ExecutorName: req.ExecutorName,
		Data:         req.Data,
		Metadata:     req.Metadata,
		Status:       models.TaskStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := s.storage.AddTask(ctx, task); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.AddTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *Service) GetTaskStatus(ctx context.Context, req *pb.GetTaskStatusRequest) (*pb.GetTaskStatusResponse, error) {
	task, err := s.storage.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if task == nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}

	return &pb.GetTaskStatusResponse{
		Status: convertTaskStatus(task.Status),
		Error:  task.Error,
	}, nil
}

// Executor Management
func (s *Service) RegisterExecutor(ctx context.Context, req *pb.RegisterExecutorRequest) (*pb.RegisterExecutorResponse, error) {
	// In a real implementation, we would track active executors
	// For now, we just verify the executor exists
	executor, err := s.storage.GetExecutor(ctx, req.ExecutorName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if executor == nil {
		return nil, status.Error(codes.NotFound, "executor not found")
	}
	if !executor.Enabled {
		return nil, status.Error(codes.FailedPrecondition, "executor is disabled")
	}

	return &pb.RegisterExecutorResponse{
		Success: true,
	}, nil
}

func (s *Service) GetNextTask(ctx context.Context, req *pb.GetNextTaskRequest) (*pb.GetNextTaskResponse, error) {
	executor, err := s.storage.GetExecutor(ctx, req.ExecutorName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if executor == nil || !executor.Enabled {
		return nil, status.Error(codes.NotFound, "executor not found or disabled")
	}
	task, err := s.storage.GetNextTask(ctx, req.ExecutorName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if task == nil {
		return nil, status.Error(codes.NotFound, "no tasks available")
	}
	return &pb.GetNextTaskResponse{
		Task: convertTaskToProto(task),
	}, nil
}

func (s *Service) UpdateTaskStatus(ctx context.Context, req *pb.UpdateTaskStatusRequest) (*pb.UpdateTaskStatusResponse, error) {
	task, err := s.storage.GetTask(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if task == nil {
		return nil, status.Error(codes.NotFound, "task not found")
	}
	executor, err := s.storage.GetExecutor(ctx, task.ExecutorName)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	taskStatus := convertProtoTaskStatus(req.Status)
	if taskStatus == models.TaskStatusFailed {
		if shouldRetry(executor.RetryPolicy, task.RetryCount) {
			taskStatus = models.TaskStatusPending
			task.RetryCount++
		} else if executor.DLQConfig.Enabled {
			if err := s.storage.MoveToDLQ(ctx, task); err != nil {
				return nil, status.Error(codes.Internal, err.Error())
			}
			return &pb.UpdateTaskStatusResponse{Task: convertTaskToProto(task)}, nil
		}
	}
	if err := s.storage.UpdateTaskStatus(ctx, req.Id, taskStatus, req.Error); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.UpdateTaskStatusResponse{Task: convertTaskToProto(task)}, nil
}

// Executor Configuration
func (s *Service) CreateExecutor(ctx context.Context, req *pb.CreateExecutorRequest) (*pb.CreateExecutorResponse, error) {
	if req == nil || req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "request or config is nil")
	}

	config := &models.ExecutorConfig{
		Name:      req.Config.Name,
		Enabled:   req.Config.Enabled,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// WriteConcern
	if req.Config.WriteConcern != nil {
		config.WriteConcern = models.WriteConcern{
			Level: convertProtoWriteConcernLevel(req.Config.WriteConcern.Level),
		}
	}

	// RetryPolicy
	if req.Config.RetryPolicy != nil {
		config.RetryPolicy = models.RetryPolicy{
			Type:        convertProtoRetryPolicyType(req.Config.RetryPolicy.Type),
			MaxAttempts: int(req.Config.RetryPolicy.MaxAttempts),
			Interval:    req.Config.RetryPolicy.Interval.AsDuration(),
		}
	}

	// DLQConfig
	if req.Config.DlqConfig != nil {
		config.DLQConfig = models.DLQConfig{
			Enabled:   req.Config.DlqConfig.Enabled,
			QueueName: req.Config.DlqConfig.QueueName,
		}
	}

	if err := s.storage.CreateExecutor(ctx, config); err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to create executor: %v", err))
	}

	return &pb.CreateExecutorResponse{
		Executor: convertExecutorToProto(config),
	}, nil
}

func (s *Service) UpdateExecutor(ctx context.Context, req *pb.UpdateExecutorRequest) (*pb.UpdateExecutorResponse, error) {
	if req == nil || req.Config == nil {
		return nil, status.Error(codes.InvalidArgument, "request or config is nil")
	}

	config := &models.ExecutorConfig{
		Name:    req.Config.Name,
		Enabled: req.Config.Enabled,
		WriteConcern: models.WriteConcern{
			Level: convertProtoWriteConcernLevel(req.Config.WriteConcern.Level),
		},
		RetryPolicy: models.RetryPolicy{
			Type:        convertProtoRetryPolicyType(req.Config.RetryPolicy.Type),
			MaxAttempts: int(req.Config.RetryPolicy.MaxAttempts),
			Interval:    req.Config.RetryPolicy.Interval.AsDuration(),
		},
		DLQConfig: models.DLQConfig{
			Enabled:   req.Config.DlqConfig.Enabled,
			QueueName: req.Config.DlqConfig.QueueName,
		},
		UpdatedAt: time.Now(),
	}

	if err := s.storage.UpdateExecutor(ctx, config); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.UpdateExecutorResponse{
		Executor: convertExecutorToProto(config),
	}, nil
}

func (s *Service) GetExecutor(ctx context.Context, req *pb.GetExecutorRequest) (*pb.GetExecutorResponse, error) {
	executor, err := s.storage.GetExecutor(ctx, req.Id)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	if executor == nil {
		return nil, status.Error(codes.NotFound, "executor not found")
	}

	return &pb.GetExecutorResponse{
		Executor: convertExecutorToProto(executor),
	}, nil
}

func (s *Service) ListExecutors(ctx context.Context, req *pb.ListExecutorsRequest) (*pb.ListExecutorsResponse, error) {
	executors, err := s.storage.ListExecutors(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	result := make([]*pb.Executor, len(executors))
	for i, executor := range executors {
		result[i] = convertExecutorToProto(executor)
	}
	return &pb.ListExecutorsResponse{
		Executors: result,
	}, nil
}

// DeleteExecutor deletes an executor by its name (id)
func (s *Service) DeleteExecutor(ctx context.Context, req *pb.DeleteExecutorRequest) (*pb.DeleteExecutorResponse, error) {
	if req == nil || req.Id == "" {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}
	if err := s.storage.DeleteExecutor(ctx, req.Id); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &pb.DeleteExecutorResponse{}, nil
}

// Helper functions
func convertTaskStatus(status models.TaskStatus) pb.TaskStatus {
	switch status {
	case models.TaskStatusPending:
		return pb.TaskStatus_TASK_STATUS_PENDING
	case models.TaskStatusInProgress:
		return pb.TaskStatus_TASK_STATUS_IN_PROGRESS
	case models.TaskStatusCompleted:
		return pb.TaskStatus_TASK_STATUS_COMPLETED
	case models.TaskStatusFailed:
		return pb.TaskStatus_TASK_STATUS_FAILED
	case models.TaskStatusDLQ:
		return pb.TaskStatus_TASK_STATUS_DLQ
	default:
		return pb.TaskStatus_TASK_STATUS_PENDING
	}
}

func convertProtoTaskStatus(status pb.TaskStatus) models.TaskStatus {
	switch status {
	case pb.TaskStatus_TASK_STATUS_PENDING:
		return models.TaskStatusPending
	case pb.TaskStatus_TASK_STATUS_IN_PROGRESS:
		return models.TaskStatusInProgress
	case pb.TaskStatus_TASK_STATUS_COMPLETED:
		return models.TaskStatusCompleted
	case pb.TaskStatus_TASK_STATUS_FAILED:
		return models.TaskStatusFailed
	case pb.TaskStatus_TASK_STATUS_DLQ:
		return models.TaskStatusDLQ
	default:
		return models.TaskStatusPending
	}
}

func convertExecutorConfig(config *models.ExecutorConfig) *pb.ExecutorConfig {
	return &pb.ExecutorConfig{
		Name:    config.Name,
		Enabled: config.Enabled,
		WriteConcern: &pb.WriteConcern{
			Level: convertWriteConcernLevel(config.WriteConcern.Level),
		},
		RetryPolicy: &pb.RetryPolicy{
			Type:        convertRetryPolicyType(config.RetryPolicy.Type),
			MaxAttempts: int32(config.RetryPolicy.MaxAttempts),
			Interval:    durationpb.New(config.RetryPolicy.Interval),
		},
		DlqConfig: &pb.DLQConfig{
			Enabled:   config.DLQConfig.Enabled,
			QueueName: config.DLQConfig.QueueName,
		},
	}
}

func convertWriteConcernLevel(level models.WriteConcernLevel) pb.WriteConcernLevel {
	switch level {
	case models.WriteConcernReplicaAcknowledged:
		return pb.WriteConcernLevel_WRITE_CONCERN_REPLICA_ACKNOWLEDGED
	case models.WriteConcernMajority:
		return pb.WriteConcernLevel_WRITE_CONCERN_MAJORITY
	case models.WriteConcernUnacknowledged:
		return pb.WriteConcernLevel_WRITE_CONCERN_UNACKNOWLEDGED
	case models.WriteConcernJournaled:
		return pb.WriteConcernLevel_WRITE_CONCERN_JOURNALED
	default:
		return pb.WriteConcernLevel_WRITE_CONCERN_REPLICA_ACKNOWLEDGED
	}
}

func convertProtoWriteConcernLevel(level pb.WriteConcernLevel) models.WriteConcernLevel {
	switch level {
	case pb.WriteConcernLevel_WRITE_CONCERN_REPLICA_ACKNOWLEDGED:
		return models.WriteConcernReplicaAcknowledged
	case pb.WriteConcernLevel_WRITE_CONCERN_MAJORITY:
		return models.WriteConcernMajority
	case pb.WriteConcernLevel_WRITE_CONCERN_UNACKNOWLEDGED:
		return models.WriteConcernUnacknowledged
	case pb.WriteConcernLevel_WRITE_CONCERN_JOURNALED:
		return models.WriteConcernJournaled
	default:
		return models.WriteConcernReplicaAcknowledged
	}
}

func convertRetryPolicyType(policyType models.RetryPolicyType) pb.RetryPolicyType {
	switch policyType {
	case models.RetryPolicyConstant:
		return pb.RetryPolicyType_RETRY_POLICY_CONSTANT
	case models.RetryPolicyLinear:
		return pb.RetryPolicyType_RETRY_POLICY_LINEAR
	case models.RetryPolicyExponential:
		return pb.RetryPolicyType_RETRY_POLICY_EXPONENTIAL
	default:
		return pb.RetryPolicyType_RETRY_POLICY_CONSTANT
	}
}

func convertProtoRetryPolicyType(policyType pb.RetryPolicyType) models.RetryPolicyType {
	switch policyType {
	case pb.RetryPolicyType_RETRY_POLICY_CONSTANT:
		return models.RetryPolicyConstant
	case pb.RetryPolicyType_RETRY_POLICY_LINEAR:
		return models.RetryPolicyLinear
	case pb.RetryPolicyType_RETRY_POLICY_EXPONENTIAL:
		return models.RetryPolicyExponential
	default:
		return models.RetryPolicyConstant
	}
}

func shouldRetry(policy models.RetryPolicy, retryCount int) bool {
	if policy.MaxAttempts == 0 {
		return true // Unlimited retries
	}
	return retryCount < policy.MaxAttempts
}

func convertTaskToProto(task *models.Task) *pb.Task {
	if task == nil {
		return nil
	}
	return &pb.Task{
		Id:           task.ID.Hex(),
		ExecutorName: task.ExecutorName,
		Data:         task.Data,
		Metadata:     task.Metadata,
		Status:       convertTaskStatus(task.Status),
		Error:        task.Error,
		RetryCount:   int32(task.RetryCount),
		CreatedAt:    timestamppb.New(task.CreatedAt),
		UpdatedAt:    timestamppb.New(task.UpdatedAt),
		StartedAt:    timestamppb.New(zeroOrTime(task.StartedAt)),
		CompletedAt:  timestamppb.New(zeroOrTime(task.CompletedAt)),
	}
}

func convertExecutorToProto(config *models.ExecutorConfig) *pb.Executor {
	if config == nil {
		return nil
	}
	return &pb.Executor{
		Id:      config.ID.Hex(),
		Name:    config.Name,
		Enabled: config.Enabled,
		Config: &pb.ExecutorConfig{
			Name:    config.Name,
			Enabled: config.Enabled,
			WriteConcern: &pb.WriteConcern{
				Level: convertWriteConcernLevel(config.WriteConcern.Level),
			},
			RetryPolicy: &pb.RetryPolicy{
				Type:        convertRetryPolicyType(config.RetryPolicy.Type),
				MaxAttempts: int32(config.RetryPolicy.MaxAttempts),
				Interval:    durationpb.New(config.RetryPolicy.Interval),
			},
			DlqConfig: &pb.DLQConfig{
				Enabled:   config.DLQConfig.Enabled,
				QueueName: config.DLQConfig.QueueName,
			},
		},
		CreatedAt: timestamppb.New(config.CreatedAt),
		UpdatedAt: timestamppb.New(config.UpdatedAt),
	}
}

func zeroOrTime(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func convertProtoExecutorConfig(config *models.ExecutorConfig) *pb.ExecutorConfig {
	if config == nil {
		return nil
	}

	result := &pb.ExecutorConfig{
		Name:    config.Name,
		Enabled: config.Enabled,
	}

	// WriteConcern
	result.WriteConcern = &pb.WriteConcern{
		Level: convertWriteConcernLevel(config.WriteConcern.Level),
	}

	// RetryPolicy
	result.RetryPolicy = &pb.RetryPolicy{
		Type:        convertRetryPolicyType(config.RetryPolicy.Type),
		MaxAttempts: int32(config.RetryPolicy.MaxAttempts),
		Interval:    durationpb.New(config.RetryPolicy.Interval),
	}

	// DLQConfig
	result.DlqConfig = &pb.DLQConfig{
		Enabled:   config.DLQConfig.Enabled,
		QueueName: config.DLQConfig.QueueName,
	}

	return result
}
