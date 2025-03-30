package repository

import (
	"context"
	"time"

	"github.com/kwa0x2/AutoSRT-Backend/domain"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type usageRepository struct {
	collection *mongo.Collection
}

func NewUsageRepository(db *mongo.Database, collection string) domain.UsageRepository {
	return &usageRepository{
		collection: db.Collection(collection),
	}
}

func (ur *usageRepository) Create(ctx context.Context, usage *domain.Usage) error {
	result, err := ur.collection.InsertOne(ctx, usage)
	if err != nil {
		return err
	}

	usage.ID = result.InsertedID.(bson.ObjectID)
	return nil
}

func (ur *usageRepository) FindOne(ctx context.Context, filter bson.D) (domain.Usage, error) {
	var usage domain.Usage

	if err := ur.collection.FindOne(ctx, filter).Decode(&usage); err != nil {
		return usage, err
	}

	return usage, nil
}

func (ur *usageRepository) UpdateOne(ctx context.Context, filter bson.D, update bson.D, opts *options.UpdateOneOptionsBuilder) error {
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
