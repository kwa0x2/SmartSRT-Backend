package usecase

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"github.com/kwa0x2/AutoSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userUseCase struct {
	userBaseRepository  domain.BaseRepository[*domain.User]
	usageBaseRepository domain.BaseRepository[*domain.Usage]
	srtBaseRepository   domain.BaseRepository[*domain.SRTHistory]
}

func NewUserUseCase(
	userBaseRepository domain.BaseRepository[*domain.User],
	usageBaseRepository domain.BaseRepository[*domain.Usage],
	srtBaseRepository domain.BaseRepository[*domain.SRTHistory],
) domain.UserUseCase {
	return &userUseCase{
		userBaseRepository:  userBaseRepository,
		usageBaseRepository: usageBaseRepository,
		srtBaseRepository:   srtBaseRepository,
	}
}

func (uu *userUseCase) Create(user *domain.User) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := uu.userBaseRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		now := time.Now().UTC()

		user.CreatedAt = now
		user.UpdatedAt = now
		user.Plan = types.Free

		if err = user.Validate(); err != nil {
			return nil, err
		}

		if err = uu.userBaseRepository.Create(txCtx, user); err != nil {
			return nil, err
		}

		usage := &domain.Usage{
			UserID:       user.ID,
			StartDate:    now,
			MonthlyUsage: float64(0),
			TotalUsage:   float64(0),
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		return nil, uu.usageBaseRepository.Create(txCtx, usage)
	}, txnOptions)

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return abortErr
		}
		return err
	}

	return nil
}

func (uu *userUseCase) FindOneByEmail(email string) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	return uu.userBaseRepository.FindOne(ctx, filter)
}

func (uu *userUseCase) FindOneByID(id bson.ObjectID) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	return uu.userBaseRepository.FindOne(ctx, filter)
}

func (uu *userUseCase) FindOneByEmailAndAuthType(email string, authType types.AuthType) (*domain.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{
		{Key: "email", Value: email},
		{Key: "auth_type", Value: authType},
	}

	return uu.userBaseRepository.FindOne(ctx, filter)
}

func (uu *userUseCase) IsEmailExists(email string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "email", Value: email}}
	_, err := uu.userBaseRepository.FindOne(ctx, filter)

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

	filter := bson.D{{Key: "phone_number", Value: phone}}
	_, err := uu.userBaseRepository.FindOne(ctx, filter)

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

	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newPassword}}}}
	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "auth_type", Value: types.Credentials},
	}

	return uu.userBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (uu *userUseCase) DeleteUser(userID bson.ObjectID) error {
	wc := writeconcern.Majority()
	txnOptions := options.Transaction().SetWriteConcern(wc)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	session, err := uu.userBaseRepository.GetDatabase().Client().StartSession()
	if err != nil {
		return err
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(txCtx context.Context) (interface{}, error) {
		userFilter := bson.D{{Key: "_id", Value: userID}}
		if err = uu.userBaseRepository.SoftDelete(txCtx, userFilter); err != nil {
			return nil, err
		}

		usageFilter := bson.D{{Key: "user_id", Value: userID}}
		if err = uu.usageBaseRepository.SoftDelete(txCtx, usageFilter); err != nil {
			return nil, err
		}

		srtFilter := bson.D{{Key: "user_id", Value: userID}}
		if err = uu.srtBaseRepository.SoftDelete(txCtx, srtFilter); err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	if err != nil {
		if abortErr := session.AbortTransaction(ctx); abortErr != nil {
			return abortErr
		}
		return err
	}

	return nil
}
