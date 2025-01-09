package repository

import (
	"context"
	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
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
