package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type subscriptionUseCase struct {
	subscriptionBaseRepository domain.BaseRepository[*domain.Subscription]
	userBaseRepository         domain.BaseRepository[*domain.User]
	usageBaseRepository        domain.BaseRepository[*domain.Usage]
}

func NewSubscriptionUseCase(subscriptionBaseRepository domain.BaseRepository[*domain.Subscription], userBaseRepository domain.BaseRepository[*domain.User], usageBaseRepository domain.BaseRepository[*domain.Usage]) domain.SubscriptionUseCase {
	return &subscriptionUseCase{
		subscriptionBaseRepository: subscriptionBaseRepository,
		userBaseRepository:         userBaseRepository,
		usageBaseRepository:        usageBaseRepository,
	}
}

func (sc *subscriptionUseCase) Create(subscription domain.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	filter := bson.D{{Key: "_id", Value: subscription.UserID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "role", Value: types.Pro},
	}}}

	if err := sc.userBaseRepository.UpdateOne(ctx, filter, update, nil); err != nil {
		return err
	}

	filter = bson.D{{Key: "user_id", Value: subscription.UserID}}
	update = bson.D{{Key: "$set", Value: bson.D{
		{Key: "monthly_usage", Value: 0},
	}}}

	if err := sc.usageBaseRepository.UpdateOne(ctx, filter, update, nil); err != nil {
		return err
	}

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

func (sc *subscriptionUseCase) UpdateBillingDatesByID(subscriptionID, nextBilledAt string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "next_billed_at", Value: nextBilledAt},
	}}}

	return sc.subscriptionBaseRepository.UpdateOne(ctx, filter, update, nil)
}
