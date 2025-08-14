package usecase

import (
	"context"
	"errors"
	"fmt"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/config"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type paddleUseCase struct {
	env                 *config.Env
	sdk                 *paddle.SDK
	subscriptionUseCase domain.SubscriptionUseCase
	customerUseCase     domain.CustomerUseCase
	userUseCase         domain.UserUseCase
}

func NewPaddleUseCase(env *config.Env, paddleSDK *paddle.SDK, subscriptionUseCase domain.SubscriptionUseCase, customerUseCase domain.CustomerUseCase, userUseCase domain.UserUseCase) domain.PaddleUseCase {
	return &paddleUseCase{
		env:                 env,
		sdk:                 paddleSDK,
		subscriptionUseCase: subscriptionUseCase,
		customerUseCase:     customerUseCase,
		userUseCase:         userUseCase,
	}
}

func (pu *paddleUseCase) HandleWebhook(event *domain.PaddleWebhookEvent) error {
	fmt.Println(event.Data)
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
	customData, ok := data["custom_data"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("custom_data is not a valid map")
	}

	userIDStr, ok := customData["user_id"].(string)
	if !ok {
		return fmt.Errorf("user_id is not a valid string")
	}

	userID, err := bson.ObjectIDFromHex(userIDStr)
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

	//userID, _ := bson.ObjectIDFromHex("689507eb588b269885ec80db")

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
		return err
	}

	userID, err := bson.ObjectIDFromHex(data["custom_data"].(map[string]interface{})["user_id"].(string))
	if err != nil {
		return fmt.Errorf("invalid user id format: %v", err)
	}

	if err = pu.userUseCase.UpdatePlanByID(userID, types.Free); err != nil {
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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
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
