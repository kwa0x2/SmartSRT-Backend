package domain

import (
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	CollectionContact = "contact"
)

type ContactCreateBody struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
	Message   string `json:"message"`
}

type Contact struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	FirstName string        `bson:"first_name" validate:"required"`
	LastName  string        `bson:"last_name"`
	Email     string        `bson:"email" validate:"required"`
	Message   string        `bson:"message" validate:"required"`
	CreatedAt time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
}

func (c *Contact) Validate() error {
	validate := validator.New()
	return validate.Struct(c)
}

func (c *Contact) GetCollectionName() string {
	return CollectionContact
}

func (c *Contact) SetID(id bson.ObjectID) {
	c.ID = id
}

type ContactUseCase interface {
	Create(domain *Contact) error
}
