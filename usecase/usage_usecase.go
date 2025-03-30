package usecase

import (
	"context"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type usageUseCase struct {
	usageRepository domain.UsageRepository
	userUseCase     domain.UserUseCase
}

func NewUsageUseCase(usageRepository domain.UsageRepository, userUseCase domain.UserUseCase) domain.UsageUseCase {
	return &usageUseCase{
		usageRepository: usageRepository,
		userUseCase:     userUseCase,
	}
}

func (uu *usageUseCase) Create(usage *domain.Usage) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()
	usage.CreatedAt = now
	usage.UpdatedAt = now
	usage.StartDate = now
	usage.MonthlyUsage = float64(0)
	usage.TotalUsage = float64(0)

	if err := usage.Validate(); err != nil {
		return err
	}

	return uu.usageRepository.Create(ctx, usage)
}

func (uu *usageUseCase) FindOneByUserID(userID bson.ObjectID) (domain.Usage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{"user_id", userID},
	}

	return uu.usageRepository.FindOne(ctx, filter)
}

func (uu *usageUseCase) UpdateUsage(userID bson.ObjectID, duration float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()

	usage, err := uu.FindOneByUserID(userID)
	if err != nil {
		usage = domain.Usage{
			UserID:       userID,
			StartDate:    now,
			MonthlyUsage: duration,
			TotalUsage:   duration,
		}
		if err = uu.Create(&usage); err != nil {
			return err
		}
		return nil
	}

	filter := bson.D{{"_id", usage.ID}}
	update := bson.D{
		{"$inc", bson.D{{"monthly_usage", duration}, {"total_usage", duration}}},
	}

	return uu.usageRepository.UpdateOne(ctx, filter, update, nil)
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
