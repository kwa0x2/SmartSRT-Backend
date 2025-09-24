package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/getsentry/sentry-go"
	sentrygin "github.com/getsentry/sentry-go/gin"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/api/route"
	"github.com/kwa0x2/SmartSRT-Backend/bootstrap"
)

func main() {
	app := bootstrap.App()
	env := app.Env
	db := app.MongoDatabase
	dynamodb := app.DynamoDB
	resendClient := app.ResendClient
	s3Client := app.S3Client
	lambdaClient := app.LambdaClient
	paddleSDK := app.PaddleSDK
	logger := slog.Default()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		logger.Info("Graceful shutdown initiated",
			slog.String("reason", "signal_received"),
		)
		sentry.Flush(2 * time.Second)
		logger.Info("Application shutting down",
			slog.String("status", "shutdown_complete"),
		)
		os.Exit(0)
	}()

	router := gin.New()

	router.Use(sentrygin.New(sentrygin.Options{
		Repanic:         true,
		WaitForDelivery: false,
		Timeout:         5 * time.Second,
	}))

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	route.Setup(env, db, dynamodb, router, resendClient, s3Client, lambdaClient, paddleSDK)
	logger.Info("Server starting",
		slog.String("address", env.ServerAddress),
		slog.String("environment", env.AppEnv),
	)
	router.Run(env.ServerAddress)
}
