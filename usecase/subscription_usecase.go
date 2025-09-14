package usecase

import (
	"context"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
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
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := sc.subscriptionBaseRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		now := time.Now().UTC()
		subscription.CreatedAt = now
		subscription.UpdatedAt = now

		if err = sc.subscriptionBaseRepository.Create(txCtx, &subscription); err != nil {
			return nil, err
		}

		filter := bson.D{{Key: "user_id", Value: subscription.UserID}}
		update := bson.D{{Key: "$set", Value: bson.D{
			{Key: "monthly_usage", Value: 0},
		}}}

		if err = sc.usageBaseRepository.UpdateOne(txCtx, filter, update, nil); err != nil {
			return nil, err
		}

		filter = bson.D{{Key: "_id", Value: subscription.UserID}}
		update = bson.D{{Key: "$set", Value: bson.D{
			{Key: "plan", Value: types.Pro},
		}}}

		return nil, sc.userBaseRepository.UpdateOne(txCtx, filter, update, nil)
	}, txnOptions)

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return abortErr
		}
		return err
	}

	return nil
}

func (sc *subscriptionUseCase) DeleteBySubsID(subscriptionID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}

	return sc.subscriptionBaseRepository.SoftDelete(ctx, filter)
}

func (sc *subscriptionUseCase) UpdateStatusBySubsID(subscriptionID, status string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "status", Value: status},
	}}}

	return sc.subscriptionBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (sc *subscriptionUseCase) UpdateCurrentBillingPeriodBySubsID(subscriptionID string, billingPeriod domain.BillingPeriod) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "subscription_id", Value: subscriptionID}}
	update := bson.D{{Key: "$set", Value: bson.D{
		{Key: "current_billing_period", Value: billingPeriod},
	}}}

	return sc.subscriptionBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (sc *subscriptionUseCase) FindByUserID(userID bson.ObjectID) (*domain.Subscription, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userID}}
	return sc.subscriptionBaseRepository.FindOne(ctx, filter)
}
