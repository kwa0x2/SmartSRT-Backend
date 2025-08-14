package domain

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID          bson.ObjectID  `bson:"_id,omitempty"`
	Name        string         `bson:"name" validate:"required"`
	Email       string         `bson:"email" validate:"required"`
	PhoneNumber string         `bson:"phone_number" validate:"required"`
	Password    string         `bson:"password"`
	AvatarURL   string         `bson:"avatar_url"`
	Plan        types.PlanType `bson:"plan" validate:"required"`
	CustomerID  string         `bson:"customer_id,omitempty"`
	AuthType    types.AuthType `bson:"auth_type"`
	LastLogin   time.Time      `bson:"last_login"`
	CreatedAt   time.Time      `bson:"created_at"  validate:"required"`
	UpdatedAt   time.Time      `bson:"updated_at"  validate:"required"`
	DeletedAt   *time.Time     `bson:"deleted_at,omitempty"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserUseCase interface {
	Create(user *User) error
	FindOneByEmail(email string) (*User, error)
	FindOneByEmailAndAuthType(email string, authType types.AuthType) (*User, error)
	FindOneByID(id bson.ObjectID) (*User, error)
	IsEmailExists(email string) (bool, error)
	IsPhoneExists(phone string) (bool, error)
	UpdateCredentialsPasswordByID(id bson.ObjectID, password string) error
	UpdatePlanByID(id bson.ObjectID, plan types.PlanType) error
	DeleteUser(id bson.ObjectID) error
}

func (u *User) GetCollectionName() string {
	return CollectionUser
}

func (u *User) SetID(id bson.ObjectID) {
	u.ID = id
}
