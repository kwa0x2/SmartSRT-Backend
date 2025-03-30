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
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	su := repository.NewSessionRepository(dynamodb, domain.TableName)
	usr := repository.NewUsageRepository(db)

	userUseCase := usecase.NewUserUseCase(ur, nil)
	usageUseCase := usecase.NewUsageUseCase(usr, userUseCase)
	userUseCase = usecase.NewUserUseCase(ur, usageUseCase)

	ud := &delivery.UserDelivery{
		UserUseCase: userUseCase,
	}

	userRoute := group.Group("/user")
	{
		userRoute.GET("/me", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, ur)), ud.GetProfileFromSession)

		userRoute.HEAD("/exists/email/:email", ud.CheckEmailExists)
		userRoute.HEAD("/exists/phone/:phone", ud.CheckPhoneExists)
	}
}
