package route

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSRTRoute(group *gin.RouterGroup, s3Client *s3.Client, lambdaClient *lambda.Client, bucketName, lambdaFuncName string, db *mongo.Database) {

	sr := repository.NewSRTRepository(s3Client, lambdaClient, db, bucketName, lambdaFuncName, domain.CollectionSRTHistory)
	su := usecase.NewSRTUseCase(sr)
	sd := &delivery.SRTDelivery{
		SRTUseCase: su,
	}

	srtRoute := group.Group("/srt")
	{
		srtRoute.POST("", sd.ConvertFileToSRT)
		srtRoute.GET("/histories", sd.FindHistories)
	}
}
