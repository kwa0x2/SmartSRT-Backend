package route

import (
	"log/slog"
	"os"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/SmartSRT-Backend/api/middleware"
	"github.com/kwa0x2/SmartSRT-Backend/bootstrap"
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/repository"
	"github.com/kwa0x2/SmartSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSRTRoute(env *config.Env, group *gin.RouterGroup, s3Client *s3.Client, lambdaClient *lambda.Client, bucketName, lambdaFuncName string, db *mongo.Database, dynamodb *dynamodb.Client) {
	logger := slog.Default()

	su := repository.NewSessionRepository(dynamodb, domain.TableName)
	sr := repository.NewSRTRepository(s3Client, lambdaClient, db, bucketName, lambdaFuncName, domain.CollectionSRTHistory)
	seu := usecase.NewSessionUseCase(su, repository.NewBaseRepository[*domain.User](db))

	usguc := usecase.NewUsageUseCase(env, repository.NewBaseRepository[*domain.Usage](db), repository.NewBaseRepository[*domain.User](db))

	rmq, err := bootstrap.NewRabbitMQ(env)
	if err != nil {
		logger.Error("RabbitMQ connection failed for SRT route",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	sd := &delivery.SRTDelivery{
		SRTUseCase: usecase.NewSRTUseCase(sr, usguc, repository.NewBaseRepository[*domain.SRTHistory](db)),
		RabbitMQ:   rmq,
	}

	srtRoute := group.Group("/srt")
	{
		srtRoute.POST("", middleware.SessionMiddleware(seu, repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), sd.ConvertFileToSRT)
		srtRoute.GET("/histories", middleware.SessionMiddleware(seu, repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), sd.FindHistories)
	}
}
