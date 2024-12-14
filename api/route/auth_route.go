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

	group.GET("auth/google/login", ad.GoogleLogin)
	group.GET("auth/google/callback", ad.GoogleCallback)
	group.GET("auth/github/login", ad.GitHubLogin)
	group.GET("auth/github/callback", ad.GitHubCallback)
}
