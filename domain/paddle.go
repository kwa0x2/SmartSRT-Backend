package domain

import (
	"context"
	"github.com/PaddleHQ/paddle-go-sdk/v3"
)

type PaddleCheckoutRequest struct {
	PlanID      string             `json:"plan_id"`
	Email       string             `json:"email"`
	CountryCode paddle.CountryCode `json:"country_code"`
	PostalCode  string             `json:"postal_code"`
}

type PaddleCustomerRequest struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type PaddleWebhookEvent struct {
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
}

type PaddleUseCase interface {
	CreateCustomer(ctx context.Context, req *PaddleCustomerRequest) (string, error)
	FindCustomerByEmail(ctx context.Context, req *PaddleCheckoutRequest) (string, error)
	FindAddressByPostalCode(ctx context.Context, customerCode string, req *PaddleCheckoutRequest) (string, error)
	CreateAddress(ctx context.Context, customerCode string, req *PaddleCheckoutRequest) (string, error)
	CreateCheckout(ctx context.Context, req *PaddleCheckoutRequest) (string, error)
}
