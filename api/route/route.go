package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

func Setup(env *bootstrap.Env, db *mongo.Database, dynamodb *dynamodb.Client, router *gin.Engine, resendClient *resend.Client, s3Client *s3.Client, lambdaClient *lambda.Client) {
	router.NoRoute(func(c *gin.Context) {
		c.String(http.StatusNotFound, "404 made by kwa -> https://github.com/kwa0x2")
	})

	groupRouter := router.Group("/api/v1")

	NewAuthRoute(env, groupRouter, db, dynamodb, resendClient)
	NewUserRoute(groupRouter, db, dynamodb)
	NewSRTRoute(groupRouter, s3Client, lambdaClient, env.AWSS3BucketName, env.AWSLambdaFuncName, db)
}
