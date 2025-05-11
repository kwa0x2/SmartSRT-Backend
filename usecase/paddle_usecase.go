package usecase

import (
	"context"
	"strings"
	"time"

	paddle "github.com/PaddleHQ/paddle-go-sdk/v3"
	"github.com/kwa0x2/AutoSRT-Backend/bootstrap"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
)

type paddleUseCase struct {
	env *bootstrap.Env
	sdk *paddle.SDK
}

func NewPaddleUseCase(env *bootstrap.Env, paddleSDK *paddle.SDK) domain.PaddleUseCase {
	return &paddleUseCase{
		env: env,
		sdk: paddleSDK,
	}
}

func (pu *paddleUseCase) CreateCustomer(ctx context.Context, req *domain.PaddleCustomerRequest) (string, error) {
	res, err := pu.sdk.CreateCustomer(ctx, &paddle.CreateCustomerRequest{
		Email: req.Email,
		Name:  &req.Name,
	})
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func (pu *paddleUseCase) FindCustomerByEmail(ctx context.Context, req *domain.PaddleCheckoutRequest) (string, error) {
	res, err := pu.sdk.ListCustomers(ctx, &paddle.ListCustomersRequest{
		Search: paddle.PtrTo(req.Email),
	})
	if err != nil {
		return "", err
	}

	var customerID string
	err = res.IterErr(ctx, func(customer *paddle.Customer) error {
		customerID = customer.ID
		return paddle.ErrStopIteration
	})

	if err != nil {
		return "", err
	}

	return customerID, nil
}

func (pu *paddleUseCase) FindAddressByPostalCode(ctx context.Context, customerCode string, req *domain.PaddleCheckoutRequest) (string, error) {
	res, err := pu.sdk.ListAddresses(ctx, &paddle.ListAddressesRequest{
		CustomerID: customerCode,
		Search:     paddle.PtrTo(req.PostalCode),
	})
	if err != nil {
		return "", err
	}

	var addressID string
	err = res.IterErr(ctx, func(address *paddle.Address) error {
		if address.CountryCode == req.CountryCode {
			addressID = address.ID
			return paddle.ErrStopIteration
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return addressID, nil
}

func (pu *paddleUseCase) CreateAddress(ctx context.Context, customerCode string, req *domain.PaddleCheckoutRequest) (string, error) {
	res, err := pu.sdk.CreateAddress(ctx, &paddle.CreateAddressRequest{
		CustomerID:  customerCode,
		CountryCode: req.CountryCode,
		PostalCode:  paddle.PtrTo(req.PostalCode),
	})
	if err != nil {
		return "", err
	}

	return res.ID, nil
}

func (pu *paddleUseCase) CreateCheckout(ctx context.Context, req *domain.PaddleCheckoutRequest) (string, error) {
	customerCode, err := pu.FindCustomerByEmail(ctx, req)
	if err != nil {
		return "", err
	}

	req.PostalCode = strings.ToUpper(strings.ReplaceAll(req.PostalCode, " ", ""))
	req.CountryCode = paddle.CountryCode(strings.ToUpper(strings.ReplaceAll(string(req.CountryCode), " ", "")))

	addressID, err := pu.FindAddressByPostalCode(ctx, customerCode, req)
	if err != nil {
		return "", err
	}

	if addressID == "" {
		addressID, err = pu.CreateAddress(ctx, customerCode, req)
		if err != nil {
			return "", err
		}
	}

	res, err := pu.sdk.CreateTransaction(ctx, &paddle.CreateTransactionRequest{
		Items: []paddle.CreateTransactionItems{
			*paddle.NewCreateTransactionItemsTransactionItemFromCatalog(&paddle.TransactionItemFromCatalog{
				Quantity: 1,
				PriceID:  req.PlanID,
			}),
		},
		AddressID:    paddle.PtrTo(addressID),
		CustomerID:   paddle.PtrTo(customerCode),
		CurrencyCode: paddle.PtrTo(paddle.CurrencyCodeUSD),
		BillingPeriod: &paddle.TimePeriod{
			StartsAt: time.Now().Format(time.RFC3339),
			EndsAt:   time.Now().AddDate(0, 1, 0).Format(time.RFC3339),
		},
	})

	if err != nil {
		return "", err
	}

	return res.ID, nil
}

//func (pu *PaddleUseCase) HandleWebhook(ctx context.Context, event *domain.PaddleWebhookEvent) error {
//	switch event.EventType {
//	case "subscription.created":
//		return pu.handleSubscriptionCreated(ctx, event.Data)
//	case "subscription.updated":
//		return pu.handleSubscriptionUpdated(ctx, event.Data)
//	case "subscription.cancelled":
//		return pu.handleSubscriptionCancelled(ctx, event.Data)
//	default:
//		return nil
//	}
//}
//
//func (pu *PaddleUseCase) handleSubscriptionCreated(ctx context.Context, data map[string]interface{}) error {
//	subscription := &domain.PaddleSubscription{
//		ID:            data["subscription_id"].(string),
//		UserID:        data["user_id"].(string),
//		PlanID:        data["plan_id"].(string),
//		Status:        data["status"].(string),
//		NextBillingAt: time.Unix(int64(data["next_billing_at"].(float64)), 0),
//	}
//
//	return pu.paddleRepository.CreateSubscription(ctx, subscription)
//}
//
//func (pu *PaddleUseCase) handleSubscriptionUpdated(ctx context.Context, data map[string]interface{}) error {
//	subscriptionID := data["subscription_id"].(string)
//	update := bson.M{
//		"status":          data["status"].(string),
//		"next_billing_at": time.Unix(int64(data["next_billing_at"].(float64)), 0),
//	}
//
//	return pu.paddleRepository.UpdateSubscription(ctx, subscriptionID, update)
//}
//
//func (pu *PaddleUseCase) handleSubscriptionCancelled(ctx context.Context, data map[string]interface{}) error {
//	subscriptionID := data["subscription_id"].(string)
//	update := bson.M{
//		"status": "cancelled",
//	}
//
//	return pu.paddleRepository.UpdateSubscription(ctx, subscriptionID, update)
//}
