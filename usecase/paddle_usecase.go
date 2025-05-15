package usecase

import (
	paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
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
	case "customer.created":
		return pu.handleCustomerCreated(event.Data)
	default:
		return nil
	}
}

func (pu *paddleUseCase) handleSubscriptionCreated(data map[string]interface{}) error {
	subscription := domain.Subscription{
		SubscriptionID: data["id"].(string),
		Status:         data["status"].(string),
		PriceID:        data["items"].([]interface{})[0].(map[string]interface{})["price"].(map[string]interface{})["id"].(string),
		ProductID:      data["items"].([]interface{})[0].(map[string]interface{})["price"].(map[string]interface{})["product_id"].(string),
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
