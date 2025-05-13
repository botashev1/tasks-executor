package executors

import (
	"context"
	"fmt"
	"log"

	"github.com/botashev/tasks-executor/pkg/manager"
	"github.com/botashev/tasks-executor/proto/taskexecutor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// OrderProcessor реализует обработчик заказов
type OrderProcessor struct {
	manager *manager.Manager
}

// NewOrderProcessor создает новый экземпляр обработчика заказов
func NewOrderProcessor(m *manager.Manager) *OrderProcessor {
	return &OrderProcessor{
		manager: m,
	}
}

// Process обрабатывает заказ
func (p *OrderProcessor) Process(ctx context.Context, taskData []byte) error {
	var order taskexecutor.Order
	if err := proto.Unmarshal(taskData, &order); err != nil {
		return fmt.Errorf("failed to unmarshal order: %w", err)
	}

	log.Printf("Processing order %s for user %s", order.Id, order.UserId)

	// Здесь должна быть ваша бизнес-логика обработки заказа
	// Например:
	// 1. Проверка наличия товаров
	// 2. Резервирование товаров
	// 3. Создание платежа
	// 4. Обновление статуса заказа
	// 5. Отправка уведомлений

	// Пример простой обработки:
	order.Status = "processed"
	order.UpdatedAt = timestamppb.Now()

	// Сохраняем обновленный заказ
	if _, err := proto.Marshal(&order); err != nil {
		return fmt.Errorf("failed to marshal updated order: %w", err)
	}

	// Здесь можно добавить сохранение в базу данных или отправку в другую систему
	log.Printf("Order %s processed successfully", order.Id)

	return nil
}

// Register регистрирует обработчик в системе
func (p *OrderProcessor) Register() error {
	config := &manager.ExecutorConfig{
		Name:    "order_processor",
		Enabled: true,
		WriteConcern: &manager.WriteConcern{
			Level: manager.WriteConcernLevel_WRITE_CONCERN_MAJORITY,
		},
		RetryPolicy: &manager.RetryPolicy{
			Type:        manager.RetryPolicyType_RETRY_POLICY_EXPONENTIAL,
			MaxAttempts: 5,
			Interval:    "1s",
		},
		DLQConfig: &manager.DLQConfig{
			Enabled:   true,
			QueueName: "order_processor_dlq",
		},
	}

	return p.manager.RegisterExecutor(config, p.Process)
}
