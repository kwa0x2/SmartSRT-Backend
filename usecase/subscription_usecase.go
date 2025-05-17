package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type subscriptionUseCase struct {
	subscriptionBaseRepository domain.BaseRepository[*domain.Subscription]
}

func NewSubscriptionUseCase(subscriptionBaseRepository domain.BaseRepository[*domain.Subscription]) domain.SubscriptionUseCase {
	return &subscriptionUseCase{subscriptionBaseRepository: subscriptionBaseRepository}
}

func (sc *subscriptionUseCase) Create(subscription domain.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	return sc.subscriptionBaseRepository.Create(ctx, &subscription)
}

func (sc *subscriptionUseCase) Delete(subscriptionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}

	return sc.subscriptionBaseRepository.SoftDeleteMany(ctx, filter)
}

func (sc *subscriptionUseCase) UpdateStatusByID(subscriptionID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: status},
	}}}

	return sc.subscriptionBaseRepository.UpdateOne(ctx, filter, update, nil)

}
