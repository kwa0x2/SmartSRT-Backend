package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/seeder"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(env *config.Env) *mongo.Database {
	logger := slog.Default()
	maxRetries := 30
	retryDelay := 2 * time.Second

	logger.Info("MongoDB connection starting",
		slog.String("database", env.MongoDBName),
		slog.Duration("timeout", 10*time.Second),
		slog.Int("max_retries", maxRetries),
	)

	var client *mongo.Client
	var database *mongo.Database

	for i := 0; i < maxRetries; i++ {
		serverAPI := options.ServerAPI(options.ServerAPIVersion1)
		opts := options.Client().ApplyURI(env.MongoURI).SetServerAPIOptions(serverAPI).SetConnectTimeout(10 * time.Second)

		var err error
		client, err = mongo.Connect(opts)
		if err != nil {
			logger.Warn("MongoDB connection attempt failed",
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.String("error", err.Error()),
				slog.String("database", env.MongoDBName),
			)

			if i == maxRetries-1 {
				logger.Error("MongoDB connection completely failed",
					slog.Int("total_attempts", maxRetries),
					slog.String("final_error", err.Error()),
					slog.String("database", env.MongoDBName),
				)
				os.Exit(1)
			}

			time.Sleep(retryDelay)
			continue
		}

		// Test connection with ping
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Ping(ctx, nil)
		cancel()
		if err != nil {
			logger.Warn("MongoDB ping attempt failed",
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.String("error", err.Error()),
				slog.String("database", env.MongoDBName),
			)

			if i == maxRetries-1 {
				logger.Error("MongoDB ping completely failed",
					slog.Int("total_attempts", maxRetries),
					slog.String("final_error", err.Error()),
					slog.String("database", env.MongoDBName),
				)
				os.Exit(1)
			}

			client.Disconnect(context.Background())
			time.Sleep(retryDelay)
			continue
		}

		// Test admin ping
		ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
		err = client.Database("admin").RunCommand(ctx, bson.D{{Key: "ping", Value: "1"}}).Err()
		cancel()
		if err != nil {
			logger.Warn("MongoDB admin ping attempt failed",
				slog.Int("attempt", i+1),
				slog.Int("max_retries", maxRetries),
				slog.String("error", err.Error()),
				slog.String("database", env.MongoDBName),
			)

			if i == maxRetries-1 {
				logger.Error("MongoDB admin ping completely failed",
					slog.Int("total_attempts", maxRetries),
					slog.String("final_error", err.Error()),
					slog.String("database", env.MongoDBName),
				)
				os.Exit(1)
			}

			client.Disconnect(context.Background())
			time.Sleep(retryDelay)
			continue
		}

		// Success
		database = client.Database(env.MongoDBName)
		logger.Info("MongoDB connection successful",
			slog.Int("attempt", i+1),
			slog.String("database", env.MongoDBName),
			slog.String("status", "connected"),
		)
		break
	}

	// Run seeder after successful connection
	s := seeder.NewSeeder(database)
	if err := s.SeedDatabase(); err != nil {
		logger.Warn("Database seeding failed",
			slog.String("error", err.Error()),
		)
	} else {
		logger.Info("Database seeding completed")
	}

	return database
}
