package route

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
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

func NewPaddleRoutes(env *config.Env, group *gin.RouterGroup, paddleSDK *paddle.SDK, db *mongo.Database, dynamodb *dynamodb.Client) {
	su := usecase.NewSubscriptionUseCase(repository.NewBaseRepository[*domain.Subscription](db), repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db))
	sr := repository.NewSessionRepository(dynamodb, domain.TableName)
	uu := usecase.NewUserUseCase(env, repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), repository.NewBaseRepository[*domain.SRTHistory](db), nil)
	pd := &delivery.PaddleDelivery{
		PaddleUseCase: usecase.NewPaddleUseCase(env, paddleSDK, su, uu),
	}

	paddleGroup := group.Group("/paddle")
	{
		paddleGroup.POST("/webhook", middleware.PaddleWebhookVerifier(env.PaddleWebhookSecretKey), pd.HandleWebhook)
		paddleGroup.GET("/customer-portal", middleware.SessionMiddleware(usecase.NewSessionUseCase(sr, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), pd.CreateCustomerPortalSessionByEmail)
		paddleGroup.GET("/price/:priceID", pd.GetPriceByID)
	}
}
