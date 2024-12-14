package route

import (
	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/api/http/delivery"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
)

func NewAuthRoute(env *bootstrap.Env, group *gin.RouterGroup) {
	ad := &delivery.AuthDelivery{
		Env: env,
	}

	group.GET("auth/login/google", ad.GoogleLogin)
	group.GET("auth/login/google/callback", ad.GoogleCallback)
}
