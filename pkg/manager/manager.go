package manager

import (
	"context"
	"time"

	pb "github.com/botashev/tasks-executor/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// Manager представляет клиент для взаимодействия с менеджером задач
type Manager struct {
	client pb.TaskExecutorManagerClient
}

// NewManager создает новый экземпляр менеджера
func NewManager(address string) (*Manager, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTaskExecutorManagerClient(conn)
	return &Manager{
		client: client,
	}, nil
}

// RegisterExecutor регистрирует обработчик в системе
func (m *Manager) RegisterExecutor(executorName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.RegisterExecutor(ctx, &pb.RegisterExecutorRequest{
		ExecutorName: executorName,
	})
	return err
}

// GetNextTask получает следующую задачу для обработки
func (m *Manager) GetNextTask(executorName string) (*pb.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := m.client.GetNextTask(ctx, &pb.GetNextTaskRequest{
		ExecutorName: executorName,
	})
	if err != nil {
		return nil, err
	}
	return resp.Task, nil
}

// UpdateTaskStatus обновляет статус задачи
func (m *Manager) UpdateTaskStatus(taskID string, status pb.TaskStatus, errorMsg string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := m.client.UpdateTaskStatus(ctx, &pb.UpdateTaskStatusRequest{
		Id:     taskID,
		Status: status,
		Error:  errorMsg,
	})
	return err
}
