package bootstrap

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"log"
	"time"
)

func ConnectMongoDB(env *Env) *mongo.Database {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(env.MongoURI).SetServerAPIOptions(serverAPI).SetConnectTimeout(10 * time.Second)

	client, err := mongo.Connect(opts)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", "1"}}).Err()
	if err != nil {
		log.Fatal(err)
	}

	database := client.Database(env.MongoDBName)
	log.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return database
}
