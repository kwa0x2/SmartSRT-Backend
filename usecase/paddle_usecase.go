package usecase

import (
	"context"
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
	customerUseCase     domain.CustomerUseCase
}

func NewPaddleUseCase(env *bootstrap.Env, paddleSDK *paddle.SDK, subscriptionUseCase domain.SubscriptionUseCase, customerUseCase domain.CustomerUseCase) domain.PaddleUseCase {
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
		fmt.Println(event.Data)
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

	item, ok := items[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item format")
	}

	price, ok := item["price"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid price format")
	}

	previouslyBilledAt := item["updated_at"].(string)
	if previouslyBilledAtValue, exists := item["previously_billed_at"]; exists && previouslyBilledAtValue != nil {
		previouslyBilledAt = previouslyBilledAtValue.(string)
	}

	subscription := domain.Subscription{
		SubscriptionID:     data["id"].(string),
		UserID:             userID,
		Status:             data["status"].(string),
		PriceID:            price["id"].(string),
		ProductID:          price["product_id"].(string),
		NextBilledAt:       data["next_billed_at"].(string),
		PreviouslyBilledAt: previouslyBilledAt,
		CustomerID:         data["customer_id"].(string),
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
	if err := pu.subscriptionUseCase.UpdateStatusBySubsID(data["id"].(string), data["status"].(string)); err != nil {
		fmt.Println(err)
		return err
	}

	return pu.subscriptionUseCase.DeleteBySubsID(data["id"].(string))
}

func (pu *paddleUseCase) handleSubscriptionUpdated(data map[string]interface{}) error {
	return pu.subscriptionUseCase.UpdateBillingDatesBySubsID(
		data["id"].(string),
		data["next_billed_at"].(string),
	)
}

func (pu *paddleUseCase) CreateCustomerPortalSessionByEmail(email string) (*paddle.CustomerPortalSession, error) {
	customer, err := pu.customerUseCase.FindByEmail(email)
	if err != nil {
		return nil, err
	}

	req := &paddle.CreateCustomerPortalSessionRequest{
		CustomerID: customer.CustomerID,
	}

	session, err := pu.sdk.CreateCustomerPortalSession(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (pu *paddleUseCase) CancelSubscription(userID bson.ObjectID) error {
	subscription, err := pu.subscriptionUseCase.FindByUserID(userID)
	if err != nil {
		return err
	}

	if err = pu.customerUseCase.DeleteByCustomerID(subscription.CustomerID); err != nil {
		return err
	}

	effectiveFrom := paddle.EffectiveFromImmediately
	_, err = pu.sdk.CancelSubscription(context.Background(), &paddle.CancelSubscriptionRequest{
		SubscriptionID: subscription.SubscriptionID,
		EffectiveFrom:  &effectiveFrom,
	})

	if err != nil {
		return err
	}

	return nil
}
