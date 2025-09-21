package delivery

import (
	"log/slog"
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
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to bind JSON for Paddle webhook",
				slog.String("action", "validation_paddle_webhook"),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Invalid request body. Please check your input."))
		return
	}

	if err := pd.PaddleUseCase.HandleWebhook(&event); err != nil {
		if !utils.IsNormalBusinessError(err) {
			slog.Error("Failed to process Paddle webhook",
				slog.String("action", "paddle_webhook_processing"),
				slog.String("error", err.Error()))
		}
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse(err.Error()))
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
		slog.Error("Failed to create customer portal session",
			slog.String("action", "customer_portal_session_creation"),
			slog.String("email", userData.Email),
			slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to create customer portal session"))
		return
	}

	ctx.JSON(http.StatusOK, session)
}

func (pd *PaddleDelivery) GetPriceByID(ctx *gin.Context) {
	priceID := ctx.Param("priceID")
	if priceID == "" {
		ctx.JSON(http.StatusBadRequest, utils.NewMessageResponse("Price ID is required"))
		return
	}

	price, err := pd.PaddleUseCase.GetPriceByID(priceID)
	if err != nil {
		slog.Error("Failed to get price by ID",
			slog.String("action", "get_price_by_id"),
			slog.String("price_id", priceID),
			slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to get price information"))
		return
	}
	
	convertedAmount, err := utils.ConvertCentsToDollars(price.UnitPrice.Amount)
	if err != nil {
		slog.Error("Failed to parse price amount",
			slog.String("action", "parse_price_amount"),
			slog.String("price_id", priceID),
			slog.String("amount", price.UnitPrice.Amount),
			slog.String("error", err.Error()))
		ctx.JSON(http.StatusInternalServerError, utils.NewMessageResponse("Failed to parse price amount"))
		return
	}
	price.UnitPrice.Amount = convertedAmount

	ctx.JSON(http.StatusOK, price)
}
