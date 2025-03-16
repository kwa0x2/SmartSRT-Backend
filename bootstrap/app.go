package bootstrap

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Env           *Env
	MongoDatabase *mongo.Database
	DynamoDB      *dynamodb.Client
	ResendClient  *resend.Client
	S3Client      *s3.Client
	LambdaClient  *lambda.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)
	app.DynamoDB = InitDynamoDB(app.Env)
	app.ResendClient = resend.NewClient(app.Env.ResendApiKey)
	app.S3Client = s3.NewFromConfig(AWSConfig(app.Env))
	app.LambdaClient = lambda.NewFromConfig(AWSConfig(app.Env))

	return app
}
