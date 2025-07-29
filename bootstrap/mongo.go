package bootstrap

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func ConnectMongoDB(env *config.Env) *mongo.Database {
	logger := slog.Default()

	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(env.MongoURI).SetServerAPIOptions(serverAPI).SetConnectTimeout(10 * time.Second)

	logger.Info("MongoDB connection starting",
		slog.String("database", env.MongoDBName),
		slog.Duration("timeout", 10*time.Second),
	)

	client, err := mongo.Connect(opts)
	if err != nil {
		logger.Error("MongoDB connection failed",
			slog.String("error", err.Error()),
			slog.String("database", env.MongoDBName),
		)
		os.Exit(1)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		logger.Error("MongoDB ping failed",
			slog.String("error", err.Error()),
			slog.String("database", env.MongoDBName),
		)
		os.Exit(1)
	}

	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{Key: "ping", Value: "1"}}).Err()
	if err != nil {
		logger.Error("MongoDB admin ping failed",
			slog.String("error", err.Error()),
			slog.String("database", env.MongoDBName),
		)
		os.Exit(1)
	}

	database := client.Database(env.MongoDBName)

	logger.Info("MongoDB connection successful",
		slog.String("database", env.MongoDBName),
		slog.String("status", "connected"),
	)

	return database
}
