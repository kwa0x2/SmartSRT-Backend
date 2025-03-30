package domain

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	CollectionUsage = "usage"
)

type Usage struct {
	ID        bson.ObjectID `bson:"_id,omitempty"`
	UserID    bson.ObjectID `bson:"user_id" validate:"required"`
	Month     time.Time     `bson:"month" validate:"required"`
	TotalTime float64       `bson:"total_time"`
	CreatedAt time.Time     `bson:"created_at" validate:"required"`
	UpdatedAt time.Time     `bson:"updated_at" validate:"required"`
}

func (u *Usage) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

type UsageRepository interface {
	Create(ctx context.Context, usage *Usage) error
	FindOne(ctx context.Context, filter bson.D) (Usage, error)
	UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error
}

type UsageUseCase interface {
	Create(usage *Usage) error
	FindOneByUserID(userID bson.ObjectID) (Usage, error)
	UpdateUsage(userID bson.ObjectID, duration float64) error
	CheckUsageLimit(userID bson.ObjectID, duration float64) (bool, error)
}
