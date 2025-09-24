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

func NewUserRoute(env *config.Env, group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client) {
	sr := repository.NewSessionRepository(dynamodb, domain.TableName)

	ud := &delivery.UserDelivery{
		UserUseCase:        usecase.NewUserUseCase(env, repository.NewBaseRepository[*domain.User](db), nil, nil, nil),
		UserBaseRepository: repository.NewBaseRepository[*domain.User](db),
	}

	userRoute := group.Group("/user")
	{
		userRoute.GET("/me", middleware.SessionMiddleware(usecase.NewSessionUseCase(sr, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), env), ud.GetProfileFromSession)

		userRoute.HEAD("/exists/email/:email", ud.CheckEmailExists)
		userRoute.HEAD("/exists/phone/:phone", ud.CheckPhoneExists)
	}
}
