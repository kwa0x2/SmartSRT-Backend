package route

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func SetupPaddleRoutes(env *bootstrap.Env, group *gin.RouterGroup, paddleSDK *paddle.SDK, db *mongo.Database) {
	su := usecase.NewSubscriptionUseCase(repository.NewBaseRepository[*domain.Subscription](db))
	cu := usecase.NewCustomerUseCase(repository.NewBaseRepository[*domain.Customer](db))

	pd := &delivery.PaddleDelivery{
		PaddleUseCase: usecase.NewPaddleUseCase(env, paddleSDK, su, cu),
	}

	paddleGroup := group.Group("/paddle")
	{
		paddleGroup.POST("/webhook", pd.HandleWebhook)
	}
}
