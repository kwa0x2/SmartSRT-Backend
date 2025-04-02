package domain

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Entity interface {
	GetCollectionName() string
	SetID(id bson.ObjectID)
}

type BaseRepository[T Entity] interface {
	Create(ctx context.Context, entity T) error
	FindOne(ctx context.Context, filter bson.D) (T, error)
	FindOneByID(ctx context.Context, filter bson.D, userID bson.ObjectID) (T, error)
	Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]T, error)
	UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error
	SoftDeleteMany(ctx context.Context, filter bson.D) error
	GetDatabase() *mongo.Database
}
