package domain

import (
	"github.com/PaddleHQ/paddle-go-sdk/v3"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type PaddleWebhookEvent struct {
	EventType string                 `json:"event_type"`
	Data      map[string]interface{} `json:"data"`
}

type PaddleUseCase interface {
	HandleWebhook(event *PaddleWebhookEvent) error
	CreateCustomerPortalSessionByEmail(email string) (*paddle.CustomerPortalSession, error)
	CancelSubscription(userID bson.ObjectID) error
	GetCustomerIDByEmail(email string) (string, error)
}
