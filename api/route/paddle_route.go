package route

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
)

func SetupPaddleRoutes(env *bootstrap.Env, group *gin.RouterGroup, paddleSDK *paddle.SDK) {
	//su := repository.NewSessionRepository(dynamodb, domain.TableName)

	pd := &delivery.PaddleDelivery{
		PaddleUseCase: usecase.NewPaddleUseCase(env, paddleSDK),
	}

	paddleGroup := group.Group("/paddle")
	{
		//paddleGroup.POST("/checkout", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db)), pd.CreateCheckout)
		paddleGroup.POST("/checkout", pd.CreateCheckout)
		paddleGroup.POST("/customer", pd.CreateCustomer)
		//paddleGroup.POST("/webhook", pd.HandleWebhook)
	}
}
