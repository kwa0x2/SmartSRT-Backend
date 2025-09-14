package usecase

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"
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
	userUseCase         domain.UserUseCase
}

func NewPaddleUseCase(env *config.Env, paddleSDK *paddle.SDK, subscriptionUseCase domain.SubscriptionUseCase, userUseCase domain.UserUseCase) domain.PaddleUseCase {
	return &paddleUseCase{
		env:                 env,
		sdk:                 paddleSDK,
		subscriptionUseCase: subscriptionUseCase,
		userUseCase:         userUseCase,
	}
}

func (pu *paddleUseCase) HandleWebhook(event *domain.PaddleWebhookEvent) error {
	fmt.Println(event.Data)
	fmt.Println(event.EventType)

	switch event.EventType {
	case "subscription.created":
		return pu.handleSubscriptionCreated(event.Data)

	case "subscription.canceled":
		return pu.handleSubscriptionCanceled(event.Data)

	case "subscription.updated":
		if event.Data["status"].(string) == "active" {
			return pu.handleSubscriptionUpdated(event.Data)
		}

	case "subscription.past_due":
		return pu.handleSubscriptionPastDue(event.Data)
	
	case "customer.created":
		return pu.handleCustomerCreated(event.Data)
	}

	return nil
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

	billingPeriodData, ok := data["current_billing_period"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("current_billing_period is not a valid map")
	}

	startsAtStr, ok := billingPeriodData["starts_at"].(string)
	if !ok {
		return fmt.Errorf("invalid starts_at format in current_billing_period")
	}

	endsAtStr, ok := billingPeriodData["ends_at"].(string)
	if !ok {
		return fmt.Errorf("invalid ends_at format in current_billing_period")
	}

	startsAt, err := time.Parse(time.RFC3339, startsAtStr)
	if err != nil {
		return fmt.Errorf("failed to parse starts_at: %v", err)
	}

	endsAt, err := time.Parse(time.RFC3339, endsAtStr)
	if err != nil {
		return fmt.Errorf("failed to parse ends_at: %v", err)
	}

	items, ok := data["items"].([]interface{})
	if !ok || len(items) == 0 {
		return fmt.Errorf("invalid items format")
	}

	item, ok := items[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid item format")
	}

	product, ok := item["product"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid product format")
	}

	price, ok := item["price"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid price format")
	}

	unitPriceData, ok := price["unit_price"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid unit_price format")
	}

	amountCents := unitPriceData["amount"].(string)
	amountInt, err := strconv.Atoi(amountCents)
	if err != nil {
		return fmt.Errorf("failed to parse amount: %v", err)
	}
	amount := fmt.Sprintf("%.2f", float64(amountInt)/100)

	firstBilledAt, err := time.Parse(time.RFC3339, data["first_billed_at"].(string))
	if err != nil {
		return fmt.Errorf("failed to parse first_billed_at: %v", err)
	}

	subscription := domain.Subscription{
		SubscriptionID:       data["id"].(string),
		UserID:               userID,
		Status:               data["status"].(string),
		PriceID:              price["id"].(string),
		UnitPrice: domain.UnitPrice{
			Amount:       amount,
			CurrencyCode: unitPriceData["currency_code"].(string),
		},
		ProductID:            product["id"].(string),
		ProductName:          product["name"].(string),
		FirstBilledAt:        firstBilledAt.UTC(),
		CurrentBillingPeriod: domain.BillingPeriod{
			StartsAt: startsAt,
			EndsAt:   endsAt,
		},
			//subscriptionend date when the subscription is canceled
		// discount
		CustomerID:           data["customer_id"].(string),
	}
	
	if err := pu.subscriptionUseCase.Create(subscription); err != nil {
		return err
	}

	return pu.userUseCase.UpdatePlanAndUsageLimitByID(userID, types.Pro)
}

func (pu *paddleUseCase) handleCustomerCreated(data map[string]interface{}) error {
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

	return pu.userUseCase.UpdateCustomerIDByID(userID, data["id"].(string))
}

func (pu *paddleUseCase) handleSubscriptionCanceled(data map[string]interface{}) error {
	subscriptionID := data["id"].(string)

	if err := pu.subscriptionUseCase.UpdateStatusBySubsID(subscriptionID, data["status"].(string)); err != nil {
		return err
	}

	userID, err := bson.ObjectIDFromHex(data["custom_data"].(map[string]interface{})["user_id"].(string))
	if err != nil {
		return fmt.Errorf("invalid user id format: %v", err)
	}

	return pu.CancelSubscriptionImmediately(userID)
}

func (pu *paddleUseCase) handleSubscriptionUpdated(data map[string]interface{}) error {
	subscriptionID := data["id"].(string)
	
	if billingPeriodData, exists := data["current_billing_period"].(map[string]interface{}); exists {
		startsAtStr, ok := billingPeriodData["starts_at"].(string)
		if !ok {
			return fmt.Errorf("invalid starts_at format in current_billing_period")
		}

		endsAtStr, ok := billingPeriodData["ends_at"].(string)
		if !ok {
			return fmt.Errorf("invalid ends_at format in current_billing_period")
		}

		startsAt, err := time.Parse(time.RFC3339, startsAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse starts_at: %v", err)
		}

		endsAt, err := time.Parse(time.RFC3339, endsAtStr)
		if err != nil {
			return fmt.Errorf("failed to parse ends_at: %v", err)
		}

		billingPeriod := domain.BillingPeriod{
			StartsAt: startsAt,
			EndsAt:   endsAt,
		}

		if err := pu.subscriptionUseCase.UpdateCurrentBillingPeriodBySubsID(subscriptionID, billingPeriod); err != nil {
			return err
		}
	}

	return nil
}

func (pu *paddleUseCase) handleSubscriptionPastDue(data map[string]interface{}) error {
	subscriptionID := data["id"].(string)

	if err := pu.subscriptionUseCase.UpdateStatusBySubsID(subscriptionID, data["status"].(string)); err != nil {
		return err
	}

	userID, err := bson.ObjectIDFromHex(data["custom_data"].(map[string]interface{})["user_id"].(string))
	if err != nil {
		return fmt.Errorf("invalid user id format: %v", err)
	}

	return pu.CancelSubscriptionImmediately(userID)
}

func (pu *paddleUseCase) CreateCustomerPortalSessionByEmail(email string) (*paddle.CustomerPortalSession, error) {
	user, err := pu.userUseCase.FindOneByEmail(email)
	if err != nil {
		return nil, err
	}

	req := &paddle.CreateCustomerPortalSessionRequest{
		CustomerID: user.CustomerID,
	}

	session, err := pu.sdk.CreateCustomerPortalSession(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return session, nil
}

func (pu *paddleUseCase) CancelSubscriptionImmediately(userID bson.ObjectID) error {
	subscription, err := pu.subscriptionUseCase.FindByUserID(userID)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil
		}
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

	if err = pu.userUseCase.UpdatePlanAndUsageLimitByID(userID, types.Free); err != nil {
		return err
	}

	return pu.subscriptionUseCase.DeleteBySubsID(subscription.SubscriptionID)
}

func (pu *paddleUseCase) GetCustomerIDByEmail(email string) (string, error) {
	req := &paddle.ListCustomersRequest{
		Email: []string{email},
	}

	customers, err := pu.sdk.ListCustomers(context.Background(), req)
	if err != nil {
		return "", err
	}

	var customerID string

	err = customers.Iter(context.Background(), func(customer *paddle.Customer) (bool, error) {
		customerID = customer.ID
		return false, nil
	})

	if err != nil {
		return "", err
	}

	return customerID, nil
}

func (pu *paddleUseCase) GetPriceByID(priceID string) (*paddle.Price, error) {
	req := &paddle.GetPriceRequest{
		PriceID:        priceID,
		IncludeProduct: true,
	}

	price, err := pu.sdk.GetPrice(context.Background(), req)
	if err != nil {
		return nil, err
	}

	return price, nil
}
