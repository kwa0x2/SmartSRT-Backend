package delivery

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"net/http"
	"sync"
)

type AuthDelivery struct {
	Env *bootstrap.Env
}

var (
	stateStore = sync.Map{}
)

func (ad *AuthDelivery) GoogleLogin(ctx *gin.Context) {
	googleConfig := bootstrap.GoogleConfig(ad.Env)
	state := uuid.New().String()
	stateStore.Store(state, state)
	url := googleConfig.AuthCodeURL(state)
	ctx.Redirect(http.StatusTemporaryRedirect, url)
}

func (ad *AuthDelivery) GoogleCallback(ctx *gin.Context) {
	code := ctx.Query("code")
	state := ctx.Query("state")

	if _, exists := stateStore.Load(state); !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewErrorResponse("Bad Request", "Invalid state parameter. Please try again"))
		return
	}

	googleConfig := bootstrap.GoogleConfig(ad.Env)

	token, err := googleConfig.Exchange(context.Background(), code)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "Code-Token Exchange Failed"))
		return
	}

	resp, respErr := http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if respErr != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "User data fetch failed"))
		return
	}
	defer resp.Body.Close()

	var userData map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&userData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewErrorResponse("Internal Server Error", "JSON Parsing Failed"))
		return
	}

	ctx.JSON(200, userData)
}
