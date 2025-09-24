package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewSubscriptionRoute(env *config.Env, group *gin.RouterGroup, dynamodb *dynamodb.Client, db *mongo.Database) {
	su := usecase.NewSubscriptionUseCase(repository.NewBaseRepository[*domain.Subscription](db), nil, nil)
	sr := repository.NewSessionRepository(dynamodb, domain.TableName)

	sd := &delivery.SubscriptonDelivery{
		SubscriptionUseCase: su,
	}
	subscriptionRoute := group.Group("/subscription")
	{
		subscriptionRoute.GET("/remaining-days", middleware.SessionMiddleware(usecase.NewSessionUseCase(sr, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), sd.GetRemainingDays)
	}
}