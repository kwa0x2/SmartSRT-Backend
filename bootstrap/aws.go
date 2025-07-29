package bootstrap

import (
	"context"
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func AWSConfig(env *config.Env) aws.Config {
	logger := slog.Default()

	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(env.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		logger.Error("AWS configuration loading failed",
			slog.String("region", env.AWSRegion),
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	logger.Info("AWS configuration loaded successfully",
		slog.String("region", env.AWSRegion),
		slog.String("status", "loaded"),
	)

	return cfg
}
