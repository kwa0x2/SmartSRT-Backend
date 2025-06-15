package bootstrap

import (
	"context"
	"log"

	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kwa0x2/AutoSRT-Backend/config"
)

func InitDynamoDB(env *config.Env) *dynamodb.Client {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(env.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)

	}
	return dynamodb.NewFromConfig(cfg)
}
