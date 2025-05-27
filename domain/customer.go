package domain

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	CollectionCustomer = "customer"
)

type Customer struct {
	ID         bson.ObjectID `bson:"_id,omitempty"`
	CustomerID string        `bson:"customer_id,omitempty"`
	Email      string        `bson:"email" validate:"required"`
	CreatedAt  time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt  time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt  *time.Time    `bson:"deleted_at,omitempty"`
}

type CustomerUseCase interface {
	Create(customer Customer) error
	FindByEmail(email string) (*Customer, error)
	DeleteByCustomerID(customerID string) error
}

func (u *Customer) GetCollectionName() string {
	return CollectionCustomer
}

func (u *Customer) SetID(id bson.ObjectID) {
	u.ID = id
}
