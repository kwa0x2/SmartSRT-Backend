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

func NewUserRoute(group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client) {
	su := repository.NewSessionRepository(dynamodb, domain.TableName)

	ud := &delivery.UserDelivery{
		UserUseCase:        usecase.NewUserUseCase(repository.NewBaseRepository[*domain.User](db), nil, nil),
		UserBaseRepository: repository.NewBaseRepository[*domain.User](db),
	}

	userRoute := group.Group("/user")
	{
		userRoute.GET("/me", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db)), ud.GetProfileFromSession)

		userRoute.HEAD("/exists/email/:email", ud.CheckEmailExists)
		userRoute.HEAD("/exists/phone/:phone", ud.CheckPhoneExists)
	}
}
