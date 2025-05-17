package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	CollectionSubscription = "subscription"
)

type Subscription struct {
	ID                 bson.ObjectID `bson:"_id,omitempty"`
	SubscriptionID     string        `bson:"subscription_id" validate:"required"`
	UserID             bson.ObjectID `bson:"user_id" validate:"required"`
	Status             string        `bson:"status" validate:"required"`
	PriceID            string        `bson:"price_id" validate:"required"`
	ProductID          string        `bson:"product_id" validate:"required"`
	NextBilledAt       string        `bson:"next_billed_at" validate:"required"`
	PreviouslyBilledAt string        `bson:"previously_billed_at" validate:"required"`
	CustomerID         string        `bson:"customer_id" validate:"required"`
	CreatedAt          time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt          time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt          *time.Time    `bson:"deleted_at,omitempty"`
}

type SubscriptionUseCase interface {
	Create(subscription Subscription) error
	Delete(subscriptionID string) error
	UpdateStatusByID(subscriptionID, status string) error
	UpdateBillingDatesByID(subscriptionID, nextBilledAt string) error
}

func (u *Subscription) GetCollectionName() string {
	return CollectionSubscription
}

func (u *Subscription) SetID(id bson.ObjectID) {
	u.ID = id
}
