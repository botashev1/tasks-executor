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
	// Получаем адрес менеджера из переменной окружения
	managerAddr := os.Getenv("MANAGER_ADDR")
	if managerAddr == "" {
		managerAddr = "localhost:50051" // значение по умолчанию
	}

	// Создаем менеджер
	m, err := manager.NewManager(managerAddr)
	if err != nil {
		log.Fatalf("Failed to create manager: %v", err)
	}

	// Создаем пример обработчика
	exampleProcessor := executors.NewExampleProcessor(m)

	// Регистрируем обработчик
	if err := exampleProcessor.Register(); err != nil {
		log.Fatalf("Failed to register example processor: %v", err)
	}

	log.Println("Example processor started successfully")

	// Запускаем обработку задач в отдельной горутине
	go func() {
		if err := exampleProcessor.Start(); err != nil {
			log.Printf("Error in task processing: %v", err)
		}
	}()

	// Ждем сигнала для завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("Shutting down...")
}
