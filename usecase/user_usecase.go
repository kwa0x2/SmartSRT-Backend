package usecase

import (
	"context"
	"errors"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"

	"github.com/kwa0x2/SmartSRT-Backend/config"
	"github.com/kwa0x2/SmartSRT-Backend/domain"
	"github.com/kwa0x2/SmartSRT-Backend/domain/types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type userUseCase struct {
	env                 *config.Env
	userBaseRepository  domain.BaseRepository[*domain.User]
	usageBaseRepository domain.BaseRepository[*domain.Usage]
	srtBaseRepository   domain.BaseRepository[*domain.SRTHistory]
	paddleUseCase       domain.PaddleUseCase
}

func NewUserUseCase(
	env *config.Env,
	userBaseRepository domain.BaseRepository[*domain.User],
	usageBaseRepository domain.BaseRepository[*domain.Usage],
	srtBaseRepository domain.BaseRepository[*domain.SRTHistory],
	paddleUseCase domain.PaddleUseCase,
) domain.UserUseCase {
	return &userUseCase{
		env:                 env,
		userBaseRepository:  userBaseRepository,
		usageBaseRepository: usageBaseRepository,
		srtBaseRepository:   srtBaseRepository,
		paddleUseCase:       paddleUseCase,
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

		customerID, getErr := uu.paddleUseCase.GetCustomerIDByEmail(user.Email)
		if getErr != nil {
			return nil, getErr
		}
		if customerID != "" {
			user.CustomerID = customerID
		}

		if err = uu.userBaseRepository.Create(txCtx, user); err != nil {
			return nil, err
		}

		usage := &domain.Usage{
			UserID:       user.ID,
			StartDate:    now,
			MonthlyUsage: float64(0),
			TotalUsage:   float64(0),
			UsageLimit:   types.GetMonthlyLimit(user.Plan, uu.env),
			CreatedAt:    now,
			UpdatedAt:    now,
		}

		if err = uu.usageBaseRepository.Create(txCtx, usage); err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
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

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "auth_type", Value: types.Credentials},
	}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "password", Value: newPassword}}}}

	return uu.userBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (uu *userUseCase) UpdatePlanByID(id bson.ObjectID, plan types.PlanType) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "plan", Value: plan}}}}

	return uu.userBaseRepository.UpdateOne(ctx, filter, update, nil)
}

func (uu *userUseCase) UpdatePlanAndUsageLimitByID(id bson.ObjectID, plan types.PlanType) error {
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
		userFilter := bson.D{{Key: "_id", Value: id}}
		userUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "plan", Value: plan}}}}
		if err := uu.userBaseRepository.UpdateOne(txCtx, userFilter, userUpdate, nil); err != nil {
			return nil, err
		}

		usageFilter := bson.D{{Key: "user_id", Value: id}}
		usageUpdate := bson.D{{Key: "$set", Value: bson.D{{Key: "usage_limit", Value: types.GetMonthlyLimit(plan, uu.env)}}}}
		if err := uu.usageBaseRepository.UpdateOne(txCtx, usageFilter, usageUpdate, nil); err != nil {
			return nil, err
		}

		return nil, nil
	}, txnOptions)

	return err
}

func (uu *userUseCase) UpdateCustomerIDByID(id bson.ObjectID, customerID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.D{{Key: "_id", Value: id}}
	update := bson.D{{Key: "$set", Value: bson.D{{Key: "customer_id", Value: customerID}}}}

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
