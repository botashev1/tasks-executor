package seed

import (
	"context"
	"log"
	"time"

	"github.com/yourusername/tasks-executor/pkg/models"
	"github.com/yourusername/tasks-executor/pkg/storage"
)

func Executors() {
	config := storage.StorageConfig{
		MongoURI:      "mongodb://localhost:27017",
		Database:      "task_executor",
		ExecutorsColl: "executors",
		TasksColl:     "tasks",
		DLQColl:       "dlq",
	}

	s, err := storage.NewMongoStorage(config)
	if err != nil {
		log.Fatalf("Failed to create storage: %v", err)
	}

	ctx := context.Background()

	executors := []*models.ExecutorConfig{
		{
			Name:    "order_processor",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyExponential,
				MaxAttempts: 3,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "order_processor_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "payment_processor",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyExponential,
				MaxAttempts: 5,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "payment_processor_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "notification_sender",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyLinear,
				MaxAttempts: 3,
				Interval:    time.Second * 30,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "notification_sender_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "inventory_updater",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 3,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "inventory_updater_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "analytics_collector",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "analytics_collector_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "price_updater",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyExponential,
				MaxAttempts: 3,
				Interval:    time.Second * 15,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "price_updater_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "customer_sync",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyLinear,
				MaxAttempts: 3,
				Interval:    time.Second * 20,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "customer_sync_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "delivery_tracker",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 3,
				Interval:    time.Second * 30,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "delivery_tracker_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "review_processor",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "review_processor_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "stock_alert",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "stock_alert_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "discount_calculator",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "discount_calculator_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "loyalty_updater",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyExponential,
				MaxAttempts: 3,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "loyalty_updater_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "catalog_sync",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyLinear,
				MaxAttempts: 3,
				Interval:    time.Second * 15,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "catalog_sync_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "return_processor",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyExponential,
				MaxAttempts: 3,
				Interval:    time.Second * 20,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "return_processor_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "fraud_checker",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernMajority,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "fraud_checker_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "supplier_sync",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyLinear,
				MaxAttempts: 3,
				Interval:    time.Second * 30,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "supplier_sync_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "report_generator",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "report_generator_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "email_campaign",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyLinear,
				MaxAttempts: 3,
				Interval:    time.Second * 60,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "email_campaign_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "search_indexer",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 10,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "search_indexer_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
		{
			Name:    "cache_invalidator",
			Enabled: true,
			WriteConcern: models.WriteConcern{
				Level: models.WriteConcernReplicaAcknowledged,
			},
			RetryPolicy: models.RetryPolicy{
				Type:        models.RetryPolicyConstant,
				MaxAttempts: 2,
				Interval:    time.Second * 5,
			},
			DLQConfig: models.DLQConfig{
				Enabled:   true,
				QueueName: "cache_invalidator_dlq",
			},
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	for _, executor := range executors {
		if err := s.CreateExecutor(ctx, executor); err != nil {
			log.Printf("Failed to create executor %s: %v", executor.Name, err)
		} else {
			log.Printf("Created executor: %s", executor.Name)
		}
	}

	log.Println("Database seeding completed")
}
