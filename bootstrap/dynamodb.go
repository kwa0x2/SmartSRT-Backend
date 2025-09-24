package bootstrap

import (
	"context"
	"log/slog"
	"os"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kwa0x2/SmartSRT-Backend/config"
)

func InitDynamoDB(env *config.Env) *dynamodb.Client {
	logger := slog.Default()

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(env.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		logger.Error("DynamoDB SDK configuration loading failed",
			slog.String("region", env.AWSRegion),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("DynamoDB client created successfully",
		slog.String("region", env.AWSRegion),
		slog.String("status", "initialized"),
	)

	return dynamodb.NewFromConfig(cfg)
}
