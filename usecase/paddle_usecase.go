package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"
	paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/domain/types"
	"github.com/kwa0x2/SmartSRT-Backend/utils"
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
	// userID, err := utils.ParseUserIDFromCustomData(data)
	// if err != nil {
	// 	return err
	// }

	startsAt, endsAt, err := utils.ParseBillingPeriod(data)
	if err != nil {
		return err
	}

	productID, productName, priceID, amount, currencyCode, err := utils.ParseProductAndPrice(data)
	if err != nil {
		return err
	}

	firstBilledAt, err := time.Parse(time.RFC3339, data["first_billed_at"].(string))
	if err != nil {
		return fmt.Errorf("failed to parse first_billed_at: %v", err)
	}

	subscription := domain.Subscription{
		SubscriptionID:       data["id"].(string),
		UserEmail:            data["email"].(string),
		Status:               data["status"].(string),
		PriceID:              priceID,
		UnitPrice: domain.UnitPrice{
			Amount:       amount,
			CurrencyCode: currencyCode,
		},
		ProductID:            productID,
		ProductName:          productName,
		FirstBilledAt:        firstBilledAt.UTC(),
		CurrentBillingPeriod: domain.BillingPeriod{
			StartsAt: startsAt,
			EndsAt:   endsAt,
		},
		CustomerID:           data["customer_id"].(string),
	}
	
	return pu.subscriptionUseCase.Create(subscription)
}

func (pu *paddleUseCase) handleCustomerCreated(data map[string]interface{}) error {
	// userID, err := utils.ParseUserIDFromCustomData(data)
	// if err != nil {
	// 	return err
	// }

	return pu.userUseCase.UpdateCustomerIDByEmail(data["email"].(string), data["id"].(string))
}

func (pu *paddleUseCase) handleSubscriptionCanceled(data map[string]interface{}) error {
	subscriptionID := data["id"].(string)

	if err := pu.subscriptionUseCase.UpdateStatusBySubsID(subscriptionID, data["status"].(string)); err != nil {
		return err
	}

	userID, err := utils.ParseUserIDFromCustomData(data)
	if err != nil {
		return err
	}

	return pu.CancelSubscriptionImmediately(userID)
}

func (pu *paddleUseCase) handleSubscriptionUpdated(data map[string]interface{}) error {
	subscriptionID := data["id"].(string)
	
	startsAt, endsAt, err := utils.ParseBillingPeriod(data)
	if err != nil {
		return err
	}
	
	// Only update if billing period data exists (non-zero time values)
	if !startsAt.IsZero() && !endsAt.IsZero() {
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

	userID, err := utils.ParseUserIDFromCustomData(data)
	if err != nil {
		return err
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

	if subscription.Status != "canceled" {
		effectiveFrom := paddle.EffectiveFromImmediately
		_, err = pu.sdk.CancelSubscription(context.Background(), &paddle.CancelSubscriptionRequest{
			SubscriptionID: subscription.SubscriptionID,
			EffectiveFrom:  &effectiveFrom,
		})

		if err != nil {
			return err
		}
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
