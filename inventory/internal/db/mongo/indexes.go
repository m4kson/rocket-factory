package mongodb

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func EnsureIndexes(ctx context.Context, col *mongo.Collection) error {
	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "part_id", Value: 1}},
			Options: options.Index().SetUnique(true).SetName("idx_part_id_unique"),
		},
	}

	if _, err := col.Indexes().CreateMany(ctx, indexes); err != nil {
		return fmt.Errorf("ensure indexes: %w", err)
	}

	return nil
}
