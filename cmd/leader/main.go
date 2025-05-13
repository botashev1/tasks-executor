package main

import (
	"context"
	"log"
	"os"
	"time"

	pb "github.com/botashev/tasks-executor/proto"
	"google.golang.org/grpc"
)

func main() {
	managerAddr := os.Getenv("MANAGER_ADDR")
	if managerAddr == "" {
		managerAddr = "localhost:50051"
	}
	conn, err := grpc.Dial(managerAddr, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("failed to connect to manager: %v", err)
	}
	defer conn.Close()
	client := pb.NewTaskExecutorManagerClient(conn)

	leaderID := "leader-1"
	// Пример: опрашиваем список обработчиков и забираем задачи
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		resp, err := client.ListExecutors(ctx, &pb.ListExecutorsRequest{})
		cancel()
		if err != nil {
			log.Printf("failed to list executors: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		for _, exec := range resp.Executors {
			if !exec.Enabled {
				continue
			}
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			taskResp, err := client.GetNextTask(ctx, &pb.GetNextTaskRequest{
				ExecutorName: exec.Name,
				LeaderId:     leaderID,
			})
			cancel()
			if err != nil {
				continue
			}
			log.Printf("Got task for executor %s: %s", exec.Name, taskResp.TaskId)
			// Здесь должен быть вызов обработчика задачи (SDK)
			// После выполнения задачи — сообщить менеджеру о статусе
			ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
			_, err = client.UpdateTaskStatus(ctx, &pb.UpdateTaskStatusRequest{
				TaskId:       taskResp.TaskId,
				Status:       pb.TaskStatus_COMPLETED, // или FAILED
				ErrorMessage: "",
			})
			cancel()
			if err != nil {
				log.Printf("failed to update task status: %v", err)
			}
		}
		time.Sleep(2 * time.Second)
	}
}
