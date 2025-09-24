package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/SmartSRT-Backend/api/middleware"
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/repository"
	"github.com/kwa0x2/SmartSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewUsageRoute(env *config.Env, group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client) {
	sr := repository.NewSessionRepository(dynamodb, domain.TableName)

	ud := delivery.UsageDelivery{
		UsageUseCase: usecase.NewUsageUseCase(env, repository.NewBaseRepository[*domain.Usage](db), nil),
	}

	usageRoute := group.Group("/usage")
	{
		usageRoute.GET("", middleware.SessionMiddleware(usecase.NewSessionUseCase(sr, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), ud.FindOne)
	}
}
