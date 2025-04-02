package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewUsageRoute(group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client) {
	sr := repository.NewSessionRepository(dynamodb, domain.TableName)

	ud := delivery.UsageDelivery{
		UsageUseCase: usecase.NewUsageUseCase(repository.NewBaseRepository[*domain.Usage](db), nil),
	}

	usageRoute := group.Group("/usage")
	{
		usageRoute.GET("", middleware.SessionMiddleware(usecase.NewSessionUseCase(sr, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db)), ud.FindOne)
	}
}
