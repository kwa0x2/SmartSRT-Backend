package usecase

import (
	"fmt"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type paddleUseCase struct {
	env                 *bootstrap.Env
	sdk                 *paddle.SDK
	subscriptionUseCase domain.SubscriptionUseCase
	customerUseCase     domain.CustomerUsaCase
}

func NewPaddleUseCase(env *bootstrap.Env, paddleSDK *paddle.SDK, subscriptionUseCase domain.SubscriptionUseCase, customerUseCase domain.CustomerUsaCase) domain.PaddleUseCase {
	return &paddleUseCase{
		env:                 env,
		sdk:                 paddleSDK,
		subscriptionUseCase: subscriptionUseCase,
		customerUseCase:     customerUseCase,
	}
}

func (pu *paddleUseCase) HandleWebhook(event *domain.PaddleWebhookEvent) error {
	switch event.EventType {
	case "subscription.created":
		return pu.handleSubscriptionCreated(event.Data)
	case "subscription.updated":
		if event.Data["status"].(string) == "canceled" {
			return pu.handleSubscriptionCanceled(event.Data)
		} else {
			return pu.handleSubscriptionUpdated(event.Data)
		}
	case "customer.created":
		return pu.handleCustomerCreated(event.Data)

	default:
		return nil
	}
}

func (pu *paddleUseCase) handleSubscriptionCreated(data map[string]interface{}) error {
	userID, err := bson.ObjectIDFromHex(data["custom_data"].(map[string]interface{})["user_id"].(string))
	if err != nil {
		return fmt.Errorf("invalid user id format: %v", err)
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		return fmt.Errorf("invalid items format")
	}

	price, ok := items[0].(map[string]interface{})["price"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid price format")
	}

	subscription := domain.Subscription{
		SubscriptionID: data["id"].(string),
		UserID:         userID,
		Status:         data["status"].(string),
		PriceID:        price["id"].(string),
		ProductID:      price["product_id"].(string),
		NextBilledAt:   data["next_billed_at"].(string),
		CustomerID:     data["customer_id"].(string),
	}

	return pu.subscriptionUseCase.Create(subscription)
}

func (pu *paddleUseCase) handleCustomerCreated(data map[string]interface{}) error {
	customer := domain.Customer{
		CustomerID: data["id"].(string),
		Email:      data["email"].(string),
	}

	return pu.customerUseCase.Create(customer)
}

func (pu *paddleUseCase) handleSubscriptionCanceled(data map[string]interface{}) error {
	if err := pu.subscriptionUseCase.UpdateStatusByID(data["id"].(string), data["status"].(string)); err != nil {
		fmt.Println(err)
		return err
	}

	return pu.subscriptionUseCase.Delete(data["id"].(string))
}

func (pu *paddleUseCase) handleSubscriptionUpdated(data map[string]interface{}) error {
	return pu.subscriptionUseCase.UpdateBillingDatesByID(
		data["id"].(string),
		data["next_billed_at"].(string),
	)
}
