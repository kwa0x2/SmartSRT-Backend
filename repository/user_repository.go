package repository

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"
)

type userRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database, collection string) domain.UserRepository {
	return &userRepository{
		collection: db.Collection(collection),
	}
}

func (ur *userRepository) Create(ctx context.Context, user *domain.User) error {
	result, err := ur.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	user.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (ur *userRepository) FindOne(ctx context.Context, filter bson.D) (domain.User, error) {
	var user domain.User

	if err := ur.collection.FindOne(ctx, filter).Decode(&user); err != nil {
		return user, err
	}

	return user, nil

}

func (ur *userRepository) UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error {
	var err error

	update = append(update, bson.E{
		Key: "$set",
		Value: bson.D{
			{Key: "updated_at", Value: time.Now().UTC()},
		},
	})

	if opts != nil {
		_, err = ur.collection.UpdateOne(ctx, filter, update, opts)
	} else {
		_, err = ur.collection.UpdateOne(ctx, filter, update)
	}

	if err != nil {
		return err
	}

	return nil
}

func (ur *userRepository) GetDatabase() *mongo.Database {
	return ur.collection.Database()
}
