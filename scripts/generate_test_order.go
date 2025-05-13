package main

import (
	"fmt"
	"os"
	"time"

	"github.com/botashev/tasks-executor/proto/taskexecutor"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func main() {
	// Создаем тестовый заказ
	order := &taskexecutor.Order{
		Id:     "order_123",
		UserId: "user_456",
		Items: []*taskexecutor.OrderItem{
			{
				ProductId: "prod_789",
				Quantity:  2,
				Price:     99.99,
			},
			{
				ProductId: "prod_101",
				Quantity:  1,
				Price:     149.99,
			},
		},
		Total:     349.97,
		Status:    "new",
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}

	// Сериализуем в protobuf
	data, err := proto.Marshal(order)
	if err != nil {
		fmt.Printf("Failed to marshal order: %v\n", err)
		os.Exit(1)
	}

	// Сохраняем в файл
	if err := os.WriteFile("examples/order.pb", data, 0644); err != nil {
		fmt.Printf("Failed to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Test order generated successfully")
}
