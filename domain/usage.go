package domain

import (
	"context"
	"time"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
)

const (
	CollectionUsage = "usage"
)

type Usage struct {
	ID           bson.ObjectID `bson:"_id,omitempty"`
	UserID       bson.ObjectID `bson:"user_id" validate:"required"`
	StartDate    time.Time     `bson:"start_date" validate:"required"` // Subscription start date, renews every 30 days
	MonthlyUsage float64       `bson:"monthly_usage"`                  // Usage duration for current period (minutes)
	TotalUsage   float64       `bson:"total_usage"`                    // Total usage duration since registration (minutes)
	CreatedAt    time.Time     `bson:"created_at" validate:"required"`
	UpdatedAt    time.Time     `bson:"updated_at" validate:"required"`
}

func (u *Usage) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *Usage) GetCollectionName() string {
	return CollectionUsage
}

func (u *Usage) SetID(id bson.ObjectID) {
	u.ID = id
}

type UsageUseCase interface {
	FindOneByUserID(userID bson.ObjectID) (*Usage, error)
	UpdateUsage(ctx context.Context, userID bson.ObjectID, duration float64) error
	CheckUsageLimit(userID bson.ObjectID, duration float64) (bool, error)
}
