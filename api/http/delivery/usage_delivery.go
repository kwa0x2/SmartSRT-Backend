package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type UsageDelivery struct {
	UsageUseCase domain.UsageUseCase
}

func (ud *UsageDelivery) FindOne(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	usageData, err := ud.UsageUseCase.FindOneByUserID(userData.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving usage data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"MonthlyUsage": usageData.MonthlyUsage,
	})
}
