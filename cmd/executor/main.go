package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/botashev/tasks-executor/pkg/executors"
	"github.com/botashev/tasks-executor/pkg/manager"
)

func main() {
	// Создаем менеджер
	m, err := manager.NewManager("localhost:50051")
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Создаем обработчик заказов
	orderProcessor := executors.NewOrderProcessor(m)

	// Регистрируем обработчик
	if err := orderProcessor.Register(); err != nil {
		log.Fatalf("Failed to register order processor: %v", err)
	}

	log.Println("Order processor started successfully")

	// Ждем сигнала для завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
