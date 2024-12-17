package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/route"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	dynamodb := app.DynamoDB

	router := gin.New()

	route.Setup(env, db, dynamodb, router)

	router.Run(env.ServerAddress)
}
