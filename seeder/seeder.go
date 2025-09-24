package seeder

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Seeder struct {
	db     *mongo.Database
	logger *slog.Logger
}

func NewSeeder(db *mongo.Database) *Seeder {
	return &Seeder{
		db:     db,
		logger: slog.Default(),
	}
}

func (s *Seeder) SeedDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	s.logger.Info("Starting database seeding process")

	if err := s.createCollections(ctx); err != nil {
		return err
	}

	if err := s.createIndexes(ctx); err != nil {
		return err
	}

	s.logger.Info("✅ Collections and indexes created successfully")
	return nil
}

func (s *Seeder) createCollections(ctx context.Context) error {
	collections := []string{"users", "usage", "subscription"}

	for _, collName := range collections {
		err := s.db.CreateCollection(ctx, collName)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) ||
				(err.Error() != "" && (err.Error() == "collection already exists" ||
					err.Error() == "Collection already exists")) {
				s.logger.Info("Collection already exists, skipping",
					slog.String("collection", collName))
				continue
			}
			s.logger.Error("Failed to create collection",
				slog.String("collection", collName),
				slog.String("error", err.Error()))
			return err
		}
		s.logger.Info("Collection created",
			slog.String("collection", collName))
	}

	return nil
}

func (s *Seeder) createIndexes(ctx context.Context) error {
	collectionIndexes := map[string][]string{
		"users":        {"email", "phone_number", "customer_id"},
		"usage":        {"user_id"},
		"subscription": {"subscription_id", "user_id"},
	}

	for collectionName, indexFields := range collectionIndexes {
		var indexes []mongo.IndexModel

		for _, field := range indexFields {
			index := mongo.IndexModel{
				Keys: bson.D{{Key: field, Value: 1}},
				Options: options.Index().
					SetUnique(true).
					SetPartialFilterExpression(bson.D{{Key: "deleted_at", Value: nil}}),
			}
			indexes = append(indexes, index)
		}

		if err := s.createIndexesForCollection(ctx, collectionName, indexes); err != nil {
			return err
		}
	}

	return nil
}

func (s *Seeder) createIndexesForCollection(ctx context.Context, collectionName string, indexes []mongo.IndexModel) error {
	collection := s.db.Collection(collectionName)

	_, err := collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		s.logger.Error("Failed to create indexes",
			slog.String("collection", collectionName),
			slog.String("error", err.Error()))
		return err
	}

	s.logger.Info("Indexes created",
		slog.String("collection", collectionName),
		slog.Int("count", len(indexes)))

	return nil
}
