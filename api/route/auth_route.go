package route

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/api/middleware"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/repository"
	"github.com/kwa0x2/AutoSRT-Backend/usecase"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewAuthRoute(env *bootstrap.Env, group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client) {
	ur := repository.NewUserRepository(db, domain.CollectionUser)
	su := repository.NewSessionRepository(dynamodb, domain.TableName)
	sr := repository.NewSinchRepository(env.SinchAppKey, env.SinchAppSecret)
	ad := &delivery.AuthDelivery{
		Env:            env,
		UserUseCase:    usecase.NewUserUseCase(ur),
		SessionUseCase: usecase.NewSessionUseCase(su),
		SinchUseCase:   usecase.NewSinchUseCase(sr),
	}

	group.GET("auth/google/login", ad.GoogleSignIn)
	group.GET("auth/google/callback", ad.GoogleCallback)
	group.GET("auth/github/login", ad.GitHubSignIn)
	group.GET("auth/github/callback", ad.GitHubCallback)
	group.POST("auth/credentials/signup", ad.CredentialsSignUp)
	group.POST("auth/credentials/signin", ad.CredentialsSignIn)
	group.GET("auth/signout", ad.SignOut)
	group.POST("auth/sinch/send-otp", ad.SinchSendOTP)
	group.POST("auth/sinch/verify-otp", ad.SinchVerifyOTP)

	group.GET("auth/protected", middleware.SessionMiddleware(usecase.NewSessionUseCase(su)), func(c *gin.Context) {
		userID, _ := c.Get("user_id")
		role, _ := c.Get("role")
		c.JSON(200, gin.H{
			"UserID": userID,
			"Role":   role,
		})
	})
}
