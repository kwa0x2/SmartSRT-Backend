package bootstrap

import (
	"log/slog"

	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Env           *config.Env
	Logger        *slog.Logger
	MongoDatabase *mongo.Database
	DynamoDB      *dynamodb.Client
	ResendClient  *resend.Client
	S3Client      *s3.Client
	LambdaClient  *lambda.Client
	PaddleSDK     *paddle.SDK
}

func App(env *config.Env) *Application {
	app := &Application{}
	app.Env = env
	app.Logger = SetupLogger(app.Env)
	app.MongoDatabase = ConnectMongoDB(app.Env)
	app.DynamoDB = InitDynamoDB(app.Env)
	app.ResendClient = resend.NewClient(app.Env.ResendApiKey)
	app.S3Client = s3.NewFromConfig(AWSConfig(app.Env))
	app.LambdaClient = lambda.NewFromConfig(AWSConfig(app.Env))
	app.PaddleSDK = CreatePaddle(app.Env)
	return app
}
