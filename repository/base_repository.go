package repository

import (
	"context"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BaseRepository[T domain.Entity] struct {
	collection *mongo.Collection
}

func NewBaseRepository[T domain.Entity](db *mongo.Database) domain.BaseRepository[T] {
	var entity T
	collectionName := entity.GetCollectionName()
	return &BaseRepository[T]{
		collection: db.Collection(collectionName),
	}
}

func (r *BaseRepository[T]) Create(ctx context.Context, entity T) error {
	result, err := r.collection.InsertOne(ctx, entity)
	if err != nil {
		return err
	}

	if id, ok := result.InsertedID.(bson.ObjectID); ok {
		entity.SetID(id)
	}

	return nil
}

func (r *BaseRepository[T]) FindOne(ctx context.Context, filter bson.D) (T, error) {
	var entity T
	if err := r.collection.FindOne(ctx, filter).Decode(&entity); err != nil {
		return entity, err
	}

	return entity, nil
}

func (r *BaseRepository[T]) Find(ctx context.Context, filter bson.D, opts *options.FindOptionsBuilder) ([]T, error) {
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var results []T

	if err = cursor.All(ctx, &results); err != nil {
		return nil, err
	}

	return results, nil
}

func (r *BaseRepository[T]) UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error {
	var err error

	update = append(update, bson.E{
		Key: "$set",
		Value: bson.D{
			{Key: "updated_at", Value: time.Now().UTC()},
		},
	})

	if opts != nil {
		_, err = r.collection.UpdateOne(ctx, filter, update, opts)
	} else {
		_, err = r.collection.UpdateOne(ctx, filter, update)
	}

	if err != nil {
		return err
	}

	return nil
}

func (r *BaseRepository[T]) GetDatabase() *mongo.Database {
	return r.collection.Database()
}
