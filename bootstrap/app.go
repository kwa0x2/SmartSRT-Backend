package bootstrap

import (
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Application struct {
	Env           *Env
	MongoDatabase *mongo.Database
}

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	app.MongoDatabase = ConnectMongoDB(app.Env)

	return app
}
