package storage

import (
	"context"
	"errors"
	"time"

	"github.com/botashev/tasks-executor/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoStorage struct {
	client        *mongo.Client
	db            *mongo.Database
	executorsColl *mongo.Collection
	tasksColl     *mongo.Collection
	dlqColl       *mongo.Collection
}

// NewMongoStorage creates a new MongoDB storage implementation
func NewMongoStorage(config StorageConfig) (Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(config.MongoURI))
	if err != nil {
		return nil, err
	}

	db := client.Database(config.Database)
	executorsColl := db.Collection(config.ExecutorsColl)
	tasksColl := db.Collection(config.TasksColl)
	dlqColl := db.Collection(config.DLQColl)

	// Create indexes
	_, err = executorsColl.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: "name", Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, err
	}

	_, err = tasksColl.Indexes().CreateMany(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "executor_name", Value: 1}, {Key: "status", Value: 1}},
		},
		{
			Keys: bson.D{{Key: "created_at", Value: 1}},
		},
	})
	if err != nil {
		return nil, err
	}

	return &mongoStorage{
		client:        client,
		db:            db,
		executorsColl: executorsColl,
		tasksColl:     tasksColl,
		dlqColl:       dlqColl,
	}, nil
}

// Executor operations
func (s *mongoStorage) CreateExecutor(ctx context.Context, config *models.ExecutorConfig) error {
	_, err := s.executorsColl.InsertOne(ctx, config)
	return err
}

func (s *mongoStorage) UpdateExecutor(ctx context.Context, config *models.ExecutorConfig) error {
	filter := bson.M{"name": config.Name}
	update := bson.M{"$set": config}
	_, err := s.executorsColl.UpdateOne(ctx, filter, update)
	return err
}

func (s *mongoStorage) GetExecutor(ctx context.Context, name string) (*models.ExecutorConfig, error) {
	var config models.ExecutorConfig
	err := s.executorsColl.FindOne(ctx, bson.M{"name": name}).Decode(&config)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &config, nil
}

func (s *mongoStorage) ListExecutors(ctx context.Context) ([]*models.ExecutorConfig, error) {
	cursor, err := s.executorsColl.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var executors []*models.ExecutorConfig
	if err := cursor.All(ctx, &executors); err != nil {
		return nil, err
	}
	return executors, nil
}

func (s *mongoStorage) DeleteExecutor(ctx context.Context, name string) error {
	_, err := s.executorsColl.DeleteOne(ctx, bson.M{"name": name})
	return err
}

// Task operations
func (s *mongoStorage) AddTask(ctx context.Context, task *models.Task) error {
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.Status = models.TaskStatusPending
	task.RetryCount = 0

	result, err := s.tasksColl.InsertOne(ctx, task)
	if err != nil {
		return err
	}

	if oid, ok := result.InsertedID.(primitive.ObjectID); ok {
		task.ID = oid
	}
	return nil
}

func (s *mongoStorage) GetTask(ctx context.Context, id string) (*models.Task, error) {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	var task models.Task
	err = s.tasksColl.FindOne(ctx, bson.M{"_id": objectID}).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (s *mongoStorage) UpdateTaskStatus(ctx context.Context, id string, status models.TaskStatus, errorMsg string) error {
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	update := bson.M{
		"$set": bson.M{
			"status":     status,
			"error":      errorMsg,
			"updated_at": time.Now(),
		},
	}

	if status == models.TaskStatusInProgress {
		now := time.Now()
		update["$set"].(bson.M)["started_at"] = now
	} else if status == models.TaskStatusCompleted || status == models.TaskStatusFailed || status == models.TaskStatusDLQ {
		now := time.Now()
		update["$set"].(bson.M)["completed_at"] = now
	}

	_, err = s.tasksColl.UpdateOne(ctx, bson.M{"_id": objectID}, update)
	return err
}

func (s *mongoStorage) GetNextTask(ctx context.Context, executorName string) (*models.Task, error) {
	// Find and update atomically
	filter := bson.M{
		"executor_name": executorName,
		"status":        models.TaskStatusPending,
	}
	update := bson.M{
		"$set": bson.M{
			"status":     models.TaskStatusInProgress,
			"started_at": time.Now(),
			"updated_at": time.Now(),
		},
	}
	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)

	var task models.Task
	err := s.tasksColl.FindOneAndUpdate(ctx, filter, update, opts).Decode(&task)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, err
	}
	return &task, nil
}

func (s *mongoStorage) MoveToDLQ(ctx context.Context, task *models.Task) error {
	// First, update the task status
	if err := s.UpdateTaskStatus(ctx, task.ID.Hex(), models.TaskStatusDLQ, task.Error); err != nil {
		return err
	}

	// Then, copy to DLQ collection
	_, err := s.dlqColl.InsertOne(ctx, task)
	return err
}

func (s *mongoStorage) GetDLQTasks(ctx context.Context, executorName string) ([]*models.Task, error) {
	filter := bson.M{"executor_name": executorName}
	cursor, err := s.dlqColl.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var tasks []*models.Task
	if err := cursor.All(ctx, &tasks); err != nil {
		return nil, err
	}
	return tasks, nil
}

func (s *mongoStorage) ClearDLQ(ctx context.Context, executorName string) error {
	_, err := s.dlqColl.DeleteMany(ctx, bson.M{"executor_name": executorName})
	return err
}
