package delivery

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/utils"
)

type PaddleDelivery struct {
	PaddleUseCase domain.PaddleUseCase
}

func (pd *PaddleDelivery) CreateCheckout(ctx *gin.Context) {
	var req domain.PaddleCheckoutRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body"))
		return
	}

	transactionID, err := pd.PaddleUseCase.CreateCheckout(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"transaction_id": transactionID,
	})
}

func (pd PaddleDelivery) CreateCustomer(ctx *gin.Context) {
	var req domain.PaddleCustomerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body"))
	}

	_, err := pd.PaddleUseCase.CreateCustomer(ctx, &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
	}

	ctx.Status(http.StatusOK)

}

//func (pd *PaddleDelivery) HandleWebhook(ctx *gin.Context) {
//	var event domain.PaddleWebhookEvent
//	if err := ctx.ShouldBindJSON(&event); err != nil {
//		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid webhook payload"))
//		return
//	}
//
//	if err := pd.PaddleUseCase.HandleWebhook(ctx, &event); err != nil {
//		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to handle webhook"))
//		return
//	}
//
//	ctx.JSON(http.StatusOK, gin.H{
//		"message": "Webhook processed successfully",
//	})
//}
