package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	CollectionSubscription = "subscription"
)

type BillingPeriod struct {
	EndsAt   time.Time `bson:"ends_at" json:"ends_at"`
	StartsAt time.Time `bson:"starts_at" json:"starts_at"`
}

type UnitPrice struct {
	Amount       string `bson:"amount" json:"amount"`
	CurrencyCode string `bson:"currency_code" json:"currency_code"`
}

type Subscription struct {
	ID                   bson.ObjectID `bson:"_id,omitempty"`
	SubscriptionID       string        `bson:"subscription_id" validate:"required"`
	UserID               bson.ObjectID `bson:"user_id" validate:"required"`
	Status               string        `bson:"status" validate:"required"`
	PriceID              string        `bson:"price_id" validate:"required"`
	UnitPrice            UnitPrice     `bson:"unit_price" validate:"required"`
	ProductID            string        `bson:"product_id" validate:"required"`
	ProductName          string        `bson:"product_name" validate:"required"`
	FirstBilledAt        time.Time     `bson:"first_billed_at" validate:"required"` // UTC
	CurrentBillingPeriod BillingPeriod `bson:"current_billing_period" validate:"required"`
	CustomerID           string        `bson:"customer_id" validate:"required"`
	CreatedAt            time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt            time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt            *time.Time    `bson:"deleted_at,omitempty"`
}

type SubscriptionUseCase interface {
	Create(subscription Subscription) error
	DeleteBySubsID(subscriptionID string) error
	UpdateStatusBySubsID(subscriptionID, status string) error
	UpdateCurrentBillingPeriodBySubsID(subscriptionID string, billingPeriod BillingPeriod) error
	FindByUserID(userID bson.ObjectID) (*Subscription, error)
	GetRemainingDaysByUserID(userID bson.ObjectID) (int, error)
}

func (u *Subscription) GetCollectionName() string {
	return CollectionSubscription
}

func (u *Subscription) SetID(id bson.ObjectID) {
	u.ID = id
}
