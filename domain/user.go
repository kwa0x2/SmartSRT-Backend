package domain

import (
	"context"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

const (
	CollectionUser = "users"
)

type User struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	Name      string        `bson:"name" validate:"required"`
	Email     string        `bson:"email" validate:"required"`
	Password  string        `bson:"password"`
	AvatarURL string        `bson:"avatar_url" validate:"required"`
	CreatedAt time.Time     `bson:"created_at"  validate:"required"`
	UpdatedAt time.Time     `bson:"updated_at"  validate:"required"`
	DeletedAt *time.Time    `bson:"deleted_at,omitempty"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	FindOne(ctx context.Context, filter bson.D) (User, error)
}

type UserUseCase interface {
	Create(user *User) error
	FindOneByEmail(email string) (User, error)
}
