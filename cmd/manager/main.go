package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/yourusername/tasks-executor/pkg/manager"
	"github.com/yourusername/tasks-executor/pkg/storage"
	pb "github.com/yourusername/tasks-executor/proto"
	"google.golang.org/grpc"
)

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func main() {
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://localhost:27017"
	}
	storageConfig := storage.StorageConfig{
		MongoURI:      mongoURI,
		Database:      "task_executor",
		ExecutorsColl: "executors",
		TasksColl:     "tasks",
		DLQColl:       "dlq",
	}
	store, err := storage.NewMongoStorage(storageConfig)
	if err != nil {
		log.Fatalf("failed to connect to MongoDB: %v", err)
	}

	grpcServer := grpc.NewServer()
	service := manager.NewService(store)
	pb.RegisterTaskExecutorManagerServer(grpcServer, service)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	log.Println("Manager gRPC server listening on :50051")

	// REST gateway
	go func() {
		ctx := context.Background()
		mux := runtime.NewServeMux()
		opts := []grpc.DialOption{grpc.WithInsecure()}
		if err := pb.RegisterTaskExecutorManagerHandlerFromEndpoint(ctx, mux, ":50051", opts); err != nil {
			log.Fatalf("failed to start HTTP gateway: %v", err)
		}

		// Wrap the mux with CORS middleware
		handler := corsMiddleware(mux)

		log.Println("REST API listening on :8080")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
