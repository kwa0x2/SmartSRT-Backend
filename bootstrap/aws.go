package bootstrap

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"log"
)

func AWSConfig(env *config.Env) aws.Config {
	cfg, err := awsconfig.LoadDefaultConfig(context.TODO(),
		awsconfig.WithRegion(env.AWSRegion),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(env.AWSAccessKeyID, env.AWSSecretAccessKey, "")))
	if err != nil {
		log.Fatal("failed to load aws cfg")
	}

	return cfg
}
