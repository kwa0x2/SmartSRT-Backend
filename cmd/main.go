package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/route"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	dynamodb := app.DynamoDB
	resendClient := app.ResendClient
	s3Client := app.S3Client
	lambdaClient := app.LambdaClient

	router := gin.New()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	route.Setup(env, db, dynamodb, router, resendClient, s3Client, lambdaClient)

	router.Run(env.ServerAddress)
}
