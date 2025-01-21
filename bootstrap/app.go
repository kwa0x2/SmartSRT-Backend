package bootstrap

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Env           *Env
	MongoDatabase *mongo.Database
	DynamoDB      *dynamodb.Client
	ResendClient  *resend.Client
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)
	app.DynamoDB = InitDynamoDB(app.Env)
	app.ResendClient = resend.NewClient(app.Env.ResendApiKey)

	return app
}
