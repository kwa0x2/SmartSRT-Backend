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

	//region OAuth
	group.GET("auth/oauth/google/sign-in", ad.GoogleSignIn)
	group.GET("auth/oauth/google/callback", ad.GoogleCallback)
	group.GET("auth/oauth/github/sign-in", ad.GitHubSignIn)
	group.GET("auth/oauth/github/callback", ad.GitHubCallback)
	//endregion

	//region Credentials
	group.POST("auth/credentials/sign-in", ad.CredentialsSignIn)
	//endregion

	//region Create Account
	group.POST("auth/create", ad.VerifyOTPAndCreate)
	//endregion

	//region Sign Out
	group.GET("auth/sign-out", ad.SignOut)
	//endregion

	//region OTP
	group.POST("auth/otp/send", ad.SinchSendOTP)
	//endregion

	//region Check Auth
	group.GET("auth/check", middleware.SessionMiddleware(usecase.NewSessionUseCase(su)), ad.Check)
	//endregion

	//region Is Exists
	group.POST("auth/email-exists", ad.IsEmailExists)
	group.POST("auth/phone-exists", ad.IsPhoneExists)
	//endregion

}
