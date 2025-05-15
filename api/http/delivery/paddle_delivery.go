package delivery

import (
	"github.com/kwa0x2/AutoSRT-Backend/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

type PaddleDelivery struct {
	PaddleUseCase domain.PaddleUseCase
}

func (pd *PaddleDelivery) HandleWebhook(ctx *gin.Context) {
	var event domain.PaddleWebhookEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	if err := pd.PaddleUseCase.HandleWebhook(&event); err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to handle webhook"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Webhook processed successfully",
	})
}
