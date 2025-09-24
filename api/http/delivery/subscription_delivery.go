package delivery

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type SubscriptonDelivery struct {
	SubscriptionUseCase domain.SubscriptionUseCase
}

func (sd *SubscriptonDelivery) GetRemainingDays(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	days, err := sd.SubscriptionUseCase.GetRemainingDaysByUserID(userData.ID)
	if err != nil {
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to get remaining days",
				slog.String("action", "remaining_days_lookup"),
				slog.String("user_id", userData.ID.Hex()),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving remaining days. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"RemainingDays": days,
	})
}