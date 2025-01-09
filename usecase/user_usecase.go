package usecase

import (
	"context"
	"errors"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"time"
)

type userUseCase struct {
	userRepository domain.UserRepository
}

func NewUserUseCase(userRepository domain.UserRepository) domain.UserUseCase {
	return &userUseCase{
		userRepository: userRepository,
	}
}

func (uu *userUseCase) Create(user *domain.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user.CreatedAt = time.Now().UTC()
	user.UpdatedAt = time.Now().UTC()
	user.Role = types.Free
	if err := user.Validate(); err != nil {
		return err
	}
	return uu.userRepository.Create(ctx, user)
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

func (uu *userUseCase) FindOneByEmailAndAuthWith(email string, authWith types.AutoWithType) (domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{"email", email},
		{"auth_with", authWith},
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
