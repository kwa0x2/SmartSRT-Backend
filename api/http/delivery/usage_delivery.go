package delivery

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
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
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to lookup usage data",
				slog.String("action", "usage_data_lookup"),
				slog.String("user_id", userData.ID.Hex()),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving usage data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"MonthlyUsage": usageData.MonthlyUsage,
	})
}
