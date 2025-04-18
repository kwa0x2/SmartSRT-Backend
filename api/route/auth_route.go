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
	"github.com/resend/resend-go/v2"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func NewAuthRoute(env *bootstrap.Env, group *gin.RouterGroup, db *mongo.Database, dynamodb *dynamodb.Client, resendClient *resend.Client) {
	su := repository.NewSessionRepository(dynamodb, domain.TableName)
	sr := repository.NewSinchRepository(env.SinchAppKey, env.SinchAppSecret)
	rr := repository.NewResendRepository(resendClient)

	ad := &delivery.AuthDelivery{
		Env:            env,
		UserUseCase:    usecase.NewUserUseCase(repository.NewBaseRepository[*domain.User](db), repository.NewBaseRepository[*domain.Usage](db), repository.NewBaseRepository[*domain.SRTHistory](db)),
		SessionUseCase: usecase.NewSessionUseCase(su, repository.NewBaseRepository[*domain.User](db)),
		SinchUseCase:   usecase.NewSinchUseCase(sr),
		ResendUseCase:  usecase.NewResendUseCase(rr),
	}

	authGroup := group.Group("/auth")
	{
		authGroup.Use(middleware.LocaleMiddleware())
		authGroup.GET("/google/login", ad.GoogleLogin)
		authGroup.GET("/google/callback", ad.GoogleCallback)
		authGroup.GET("/github/login", ad.GitHubLogin)
		authGroup.GET("/github/callback", ad.GitHubCallback)

		authGroup.POST("/credentials/login", ad.CredentialsLogin)

		authGroup.POST("/register", ad.VerifyOTPAndCreate)

		authGroup.GET("/logout", ad.Logout)

		authGroup.POST("/otp/send", ad.SinchSendOTP)

		authGroup.POST("/account/password/forgot", ad.SendSetupNewPasswordEmail)
		authGroup.PUT("/account/password/reset", middleware.JWTMiddleware(), ad.UpdatePassword)
		authGroup.GET("/account/delete/request", middleware.SessionMiddleware(usecase.NewSessionUseCase(su, repository.NewBaseRepository[*domain.User](db)), repository.NewBaseRepository[*domain.User](db)), ad.SendDeleteAccountMail)
		authGroup.DELETE("/account", middleware.JWTMiddleware(), ad.DeleteAccount)

	}
}
