package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"log"
)

func InitDynamoDB(env *Env) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(env.AWSRegion),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		log.Fatalf("unable to load SDK config, %v", err)

	}
	return dynamodb.NewFromConfig(cfg)
}
