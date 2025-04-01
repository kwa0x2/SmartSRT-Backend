package usecase

import (
	"context"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type usageUseCase struct {
	userUseCase         domain.UserUseCase
	usageBaseRepository domain.BaseRepository[*domain.Usage]
}

func NewUsageUseCase(userUseCase domain.UserUseCase, usageBaseRepository domain.BaseRepository[*domain.Usage]) domain.UsageUseCase {
	return &usageUseCase{
		userUseCase:         userUseCase,
		usageBaseRepository: usageBaseRepository,
	}
}

func (uu *usageUseCase) FindOneByUserID(userID bson.ObjectID) (*domain.Usage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "user_id", Value: userID},
	}

	return uu.usageBaseRepository.FindOne(ctx, filter)
}

func (uu *usageUseCase) UpdateUsage(ctx context.Context, userID bson.ObjectID, duration float64) error {

	filter := bson.D{{Key: "user_id", Value: userID}}
	update := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "monthly_usage", Value: duration},
			{Key: "total_usage", Value: duration},
		}},
	}

	return uu.usageBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (uu *usageUseCase) CheckUsageLimit(userID bson.ObjectID, duration float64) (bool, error) {
	user, err := uu.userUseCase.FindOneByID(userID)
	if err != nil {
		return false, err
	}

	limit := types.GetMonthlyLimit(user.Role)

	usage, err := uu.FindOneByUserID(userID)
	if err != nil {
		return duration <= limit, nil
	}
	return (usage.MonthlyUsage + duration) <= limit, nil
}
