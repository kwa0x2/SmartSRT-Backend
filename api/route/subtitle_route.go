package route

import (
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSubtitleRoute(group *gin.RouterGroup, s3Client *s3.Client, lambdaClient *lambda.Client, bucketName, lambdaFuncName string, db *mongo.Database) {

	sr := repository.NewSubtitleRepository(s3Client, lambdaClient, bucketName, lambdaFuncName)
	su := usecase.NewSubtitleUseCase(sr)
	sd := &delivery.SubtitleDelivery{
		SubtitleUseCase: su,
	}

	subtitleRoute := group.Group("/subtitle")
	{
		subtitleRoute.POST("", sd.ConvertFileToSRT)
	}
}
