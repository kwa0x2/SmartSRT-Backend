package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSRTRoute(group *gin.RouterGroup, s3Client *s3.Client, lambdaClient *lambda.Client, bucketName, lambdaFuncName string, db *mongo.Database, dynamodb *dynamodb.Client) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	su := repository.NewSessionRepository(dynamodb, domain.TableName)
	sr := repository.NewSRTRepository(s3Client, lambdaClient, db, bucketName, lambdaFuncName, domain.CollectionSRTHistory)
	sru := usecase.NewSRTUseCase(sr)
	sd := &delivery.SRTDelivery{
		SRTUseCase: sru,
	}

	srtRoute := group.Group("/srt")
	{
		srtRoute.POST("", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, ur)), sd.ConvertFileToSRT)
		srtRoute.GET("/histories", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, ur)), sd.FindHistories)
	}
}
