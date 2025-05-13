package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
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
		fmt.Fprintf(os.Stderr, "failed to connect to manager: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	client := pb.NewTaskExecutorManagerClient(conn)

	cmd := flag.String("cmd", "", "command: add-executor | add-task | list-executors")
	name := flag.String("name", "", "executor name")
	configFile := flag.String("config", "", "executor config file (json)")
	taskFile := flag.String("task", "", "task data file (json)")
	flag.Parse()

	switch *cmd {
	case "add-executor":
		if *configFile == "" {
			fmt.Println("--config required")
			os.Exit(1)
		}
		f, err := os.ReadFile(*configFile)
		if err != nil {
			fmt.Println("failed to read config file:", err)
			os.Exit(1)
		}
		var cfg pb.ExecutorConfig
		if err := json.Unmarshal(f, &cfg); err != nil {
			fmt.Println("invalid config json:", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_, err = client.CreateExecutor(ctx, &pb.CreateExecutorRequest{Config: &cfg})
		if err != nil {
			fmt.Println("failed to create executor:", err)
			os.Exit(1)
		}
		fmt.Println("Executor created!")
	case "add-task":
		if *name == "" || *taskFile == "" {
			fmt.Println("--name and --task required")
			os.Exit(1)
		}
		f, err := os.ReadFile(*taskFile)
		if err != nil {
			fmt.Println("failed to read task file:", err)
			os.Exit(1)
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.AddTask(ctx, &pb.AddTaskRequest{
			ExecutorName: *name,
			Data:         f,
			Metadata:     map[string]string{},
		})
		if err != nil {
			fmt.Println("failed to add task:", err)
			os.Exit(1)
		}
		fmt.Println("Task added! ID:", resp.Task.Id)
	case "list-executors":
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		resp, err := client.ListExecutors(ctx, &pb.ListExecutorsRequest{})
		if err != nil {
			fmt.Println("failed to list executors:", err)
			os.Exit(1)
		}
		for _, exec := range resp.Executors {
			b, _ := json.MarshalIndent(exec, "", "  ")
			fmt.Println(string(b))
		}
	default:
		fmt.Println("Unknown or missing --cmd. Use: add-executor | add-task | list-executors")
		os.Exit(1)
	}
}
