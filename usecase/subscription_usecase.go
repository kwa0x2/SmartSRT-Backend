package usecase

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"time"
)

type subscriptionUsecase struct {
	subscriptionBaseRepository domain.BaseRepository[*domain.Subscription]
}

func NewSubscriptionUseCase(subscriptionBaseRepository domain.BaseRepository[*domain.Subscription]) domain.SubscriptionUseCase {
	return &subscriptionUsecase{subscriptionBaseRepository: subscriptionBaseRepository}
}

func (sc *subscriptionUsecase) Create(subscription domain.Subscription) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	subscription.CreatedAt = now
	subscription.UpdatedAt = now

	return sc.subscriptionBaseRepository.Create(ctx, &subscription)
}
