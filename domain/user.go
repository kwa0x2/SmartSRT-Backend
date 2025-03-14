package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
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
	Role        types.RoleType `bson:"role" validate:"required"`
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

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindOne(ctx context.Context, filter bson.D) (User, error)
	UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error
}

type UserUseCase interface {
	Create(user *User) error
	FindOneByEmail(email string) (User, error)
	FindOneByEmailAndAuthType(email string, authType types.AuthType) (User, error)
	FindOneByID(id bson.ObjectID) (User, error)
	IsEmailExists(email string) (bool, error)
	IsPhoneExists(phone string) (bool, error)
	UpdateCredentialsPasswordByID(id bson.ObjectID, newPassword string) error
}
