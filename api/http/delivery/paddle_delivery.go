package delivery

import (
	"net/http"

	"github.com/kwa0x2/AutoSRT-Backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

type PaddleDelivery struct {
	PaddleUseCase domain.PaddleUseCase
}

func (pd *PaddleDelivery) HandleWebhook(ctx *gin.Context) {
	var event domain.PaddleWebhookEvent
	if err := ctx.ShouldBindJSON(&event); err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "json_binding_paddle_webhook"})
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	if err := pd.PaddleUseCase.HandleWebhook(&event); err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "paddle_webhook_processing"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to handle webhook"))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "Webhook processed successfully",
	})
}

func (pd *PaddleDelivery) CreateCustomerPortalSessionByEmail(ctx *gin.Context) {
	user, exists := ctx.Get("user")
	if !exists {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("An error occurred. Please try again later or contact support."))
		return
	}

	userData := user.(*domain.User)

	session, err := pd.PaddleUseCase.CreateCustomerPortalSessionByEmail(userData.Email)
	if err != nil {
		utils.HandleErrorWithSentry(ctx, err, map[string]interface{}{"action": "create_customer_portal_session"})
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to create customer portal session"))
		return
	}

	ctx.JSON(http.StatusOK, session)
}
