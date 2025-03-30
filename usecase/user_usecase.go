package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userUseCase struct {
	userRepository domain.UserRepository
	usageUseCase   domain.UsageUseCase
}

func NewUserUseCase(userRepository domain.UserRepository, usageUseCase domain.UsageUseCase) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
		usageUseCase:   usageUseCase,
	}
}

func (uu *userUseCase) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	now := time.Now().UTC()

	user.CreatedAt = now
	user.UpdatedAt = now
	user.Role = types.Free

	if err := user.Validate(); err != nil {
		return err
	}

	if err := uu.userRepository.Create(ctx, user); err != nil {
		return err
	}

	usage := &domain.Usage{
		UserID:    user.ID,
		StartDate:     now,
		MonthlyUsage: float64(0),
		TotalUsage: float64(0),
	}

	return uu.usageUseCase.Create(usage)
}

func (uu *userUseCase) FindOneByEmail(email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"email", email}}
	result, err := uu.userRepository.FindOne(ctx, filter)
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (uu *userUseCase) FindOneByEmailAndAuthType(email string, authType types.AuthType) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{"email", email},
		{"auth_type", authType},
	}

	result, err := uu.userRepository.FindOne(ctx, filter)
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (uu *userUseCase) FindOneByID(id bson.ObjectID) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"_id", id}}
	result, err := uu.userRepository.FindOne(ctx, filter)
	if err != nil {
		return domain.User{}, err
	}
	return result, nil
}

func (uu *userUseCase) IsEmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"email", email}}
	_, err := uu.userRepository.FindOne(ctx, filter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (uu *userUseCase) IsPhoneExists(phone string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{"phone_number", phone}}
	_, err := uu.userRepository.FindOne(ctx, filter)

	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (uu *userUseCase) UpdateCredentialsPasswordByID(id bson.ObjectID, newPassword string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	update := bson.D{{"$set", bson.D{{"password", newPassword}}}}
	filter := bson.D{
		{"_id", id},
		{"auth_type", types.Credentials},
	}
	if err := uu.userRepository.UpdateOne(ctx, filter, update, nil); err != nil {
		return err
	}

	return nil
}
