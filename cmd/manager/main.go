package main

import (
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	pb "github.com/botashev/tasks-executor/proto"

	"github.com/botashev/tasks-executor/pkg/manager"
	"github.com/botashev/tasks-executor/pkg/storage"
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

	// HTTP server для фронтенда
	go func() {
		// Новый ServeMux для объединения API и статики
		mux := http.NewServeMux()
		// admin_page.html отдаётся по / и /admin_page.html
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/" || r.URL.Path == "/admin_page.html" {
				http.ServeFile(w, r, "./frontend/admin_page.html")
				return
			}
			http.NotFound(w, r)
		})
		// Остальные файлы (assets, components) отдаются по своим путям
		mux.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./frontend/assets"))))
		mux.Handle("/components/", http.StripPrefix("/components", http.FileServer(http.Dir("./frontend/components"))))

		// API routes
		api := http.NewServeMux()
		api.HandleFunc("/executors", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.Path)
			if r.Method == http.MethodGet {
				// Get all executors
				resp, err := service.ListExecutors(r.Context(), &pb.ListExecutorsRequest{})
				if err != nil {
					log.Printf("Error listing executors: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			} else if r.Method == http.MethodPost {
				// Create new executor
				var req pb.CreateExecutorRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					log.Printf("Error decoding request: %v", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				created, err := service.CreateExecutor(r.Context(), &req)
				if err != nil {
					log.Printf("Error creating executor: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(created)
			}
		})

		api.HandleFunc("/executors/", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("Received request: %s %s", r.Method, r.URL.Path)
			if r.URL.Path == "/executors/" {
				if r.Method == http.MethodGet {
					resp, err := service.ListExecutors(r.Context(), &pb.ListExecutorsRequest{})
					if err != nil {
						log.Printf("Error listing executors: %v", err)
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
					w.Header().Set("Content-Type", "application/json")
					json.NewEncoder(w).Encode(resp)
					return
				}
				// Можно добавить POST для создания, если нужно
			}

			// Новый универсальный способ извлечения id
			id := strings.TrimPrefix(r.URL.Path, "/executors/")
			id = strings.TrimSuffix(id, "/")
			if id == "" {
				http.Error(w, "Executor id required", http.StatusBadRequest)
				return
			}
			log.Printf("Processing request for executor: %s", id)

			switch r.Method {
			case http.MethodGet:
				resp, err := service.GetExecutor(r.Context(), &pb.GetExecutorRequest{Id: id})
				if err != nil {
					log.Printf("Error getting executor: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(resp)
			case http.MethodPut:
				var req pb.UpdateExecutorRequest
				if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
					log.Printf("Error decoding request: %v", err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
				req.Id = id
				updated, err := service.UpdateExecutor(r.Context(), &req)
				if err != nil {
					log.Printf("Error updating executor: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Header().Set("Content-Type", "application/json")
				json.NewEncoder(w).Encode(updated)
			case http.MethodDelete:
				_, err := service.DeleteExecutor(r.Context(), &pb.DeleteExecutorRequest{Id: id})
				if err != nil {
					log.Printf("Error deleting executor: %v", err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.WriteHeader(http.StatusNoContent)
			}
		})

		// Mount API routes with logging
		apiHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("API request received: %s %s", r.Method, r.URL.Path)
			http.StripPrefix("/api/v1", corsMiddleware(api)).ServeHTTP(w, r)
		})
		mux.Handle("/api/v1/", apiHandler)

		// Оборачиваем всё в CORS middleware
		handler := corsMiddleware(mux)

		log.Println("Frontend available at http://localhost:8080")
		if err := http.ListenAndServe(":8080", handler); err != nil {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
