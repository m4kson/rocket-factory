package e2e

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/m4kson/rocket-factory/inventory/internal/repository/model"
	inventory_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partId := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"part_id":        partId,
		"name":           "Test Rocket Engine",
		"description":    "Integration test part",
		"price":          float32(9999.99),
		"stock_quantity": int64(100),
		"category":       model.CategoryEngine,

		"dimensions": bson.M{
			"length": 2.5,
			"width":  1.2,
			"height": 1.4,
			"weight": 850.0,
		},

		"manufacturer": bson.M{
			"name":    "SpaceY",
			"country": "USA",
			"website": "https://spacey.test",
		},

		"tags": []string{
			"test",
			"integration",
			"rocket",
		},

		"metadata": bson.M{
			"is_test_data": true,
			"batch":        "integration-tests",
			"priority":     1,
		},

		"created_at": now,
		"updated_at": now,
	}

	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partId, nil
}

func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, info *inventory_v1.Part) (string, error) {
	partId := gofakeit.UUID()

	partDoc := bson.M{
		"part_id":        partId,
		"name":           info.GetName(),
		"description":    info.GetDescription(),
		"price":          info.GetPrice(),
		"stock_quantity": info.GetStockQuantity(),
		"category":       info.GetCategory().String(),

		"dimensions": bson.M{
			"length": info.GetDimensions().GetLength(),
			"width":  info.GetDimensions().GetWidth(),
			"height": info.GetDimensions().GetHeight(),
			"weight": info.GetDimensions().GetWeight(),
		},

		"manufacturer": bson.M{
			"name":    info.GetManufacturer().GetName(),
			"country": info.GetManufacturer().GetCountry(),
			"website": info.GetManufacturer().GetWebsite(),
		},

		"tags": info.GetTags(),

		"metadata": info.GetMetadata(),

		"created_at": info.GetCreatedAt().AsTime(),
		"updated_at": info.GetUpdatedAt().AsTime(),
	}

	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partId, nil
}

func (env *TestEnvironment) GetTestPartInfo() *inventory_v1.Part {
	now := time.Now()
	partID := gofakeit.UUID()
	return &inventory_v1.Part{
		Uuid:          partID,
		Name:          "Test Rocket Engine",
		Description:   "Integration test part",
		Price:         9999.99,
		StockQuantity: 100,

		Category: &inventory_v1.Category{
			Category: &inventory_v1.Category_Engine{
				Engine: "engine",
			},
		},

		Dimensions: &inventory_v1.Dimensions{
			Length: 2.5,
			Width:  1.2,
			Height: 1.4,
			Weight: 850.0,
		},

		Manufacturer: &inventory_v1.Manufacturer{
			Name:    "SpaceY",
			Country: "USA",
			Website: "https://spacey.test",
		},

		Tags: []string{
			"test",
			"integration",
			"rocket",
		},

		Metadata: map[string]*inventory_v1.Value{
			"is_test_data": {
				Value: &inventory_v1.Value_BoolValue{BoolValue: true},
			},
			"batch": {
				Value: &inventory_v1.Value_StringValue{StringValue: "integration-tests"},
			},
			"priority": {
				Value: &inventory_v1.Value_Int64Value{Int64Value: 1},
			},
		},

		CreatedAt: timestamppb.New(now),
		UpdatedAt: timestamppb.New(now),
	}
}

func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	databaseName := os.Getenv("MONGO_INITDB_DATABASE")
	if databaseName == "" {
		databaseName = "inventory"
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(inventoryCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
