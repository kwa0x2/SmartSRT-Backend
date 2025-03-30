package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type UsageDelivery struct {
	UsageUseCase domain.UsageUseCase
}

func (ud *UsageDelivery) FindOne(ctx *gin.Context) {
	sessionUserID, exists := ctx.Get("user_id")
	if !exists {
		ctx.JSON(http.StatusUnauthorized, utils.NewMessageResponse("Unauthorized. Please log in and try again."))
		return
	}

	userIDStr, ok := sessionUserID.(string)
	if !ok {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid session data. Please log in again. If the issue persists, contact support."))
		return
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	usageData, err := ud.UsageUseCase.FindOneByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("An error occurred while retrieving usage data. Please try again later or contact support."))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"MonthlyUsage": usageData.MonthlyUsage,
	})
}
