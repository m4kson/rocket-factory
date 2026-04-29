package app

import (
	"context"
	"fmt"
	"time"

	inventoryV1API "github.com/m4kson/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/m4kson/rocket-factory/inventory/internal/config"
	mongodb "github.com/m4kson/rocket-factory/inventory/internal/db/mongo"
	"github.com/m4kson/rocket-factory/inventory/internal/repository"
	inventoryRepository "github.com/m4kson/rocket-factory/inventory/internal/repository/part"
	"github.com/m4kson/rocket-factory/inventory/internal/service"
	inventoryService "github.com/m4kson/rocket-factory/inventory/internal/service/part"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API inventoryV1.InventoryServiceServer

	inventoryService    service.PartService
	inventoryRepository repository.PartRepository

	mongoDBClient *mongodb.Client
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryV1API(ctx context.Context) inventoryV1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = inventoryV1API.NewAPI(d.PartService(ctx))
	}

	return d.inventoryV1API
}

func (d *diContainer) PartService(ctx context.Context) service.PartService {
	if d.inventoryService == nil {
		d.inventoryService = inventoryService.NewPartService(d.PartRepository(ctx))
	}

	return d.inventoryService
}

func (d *diContainer) PartRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
		inventoryCol := d.MongoDBClient(ctx).Collection(config.AppConfig().Mongo.DbName())
		if err := mongodb.EnsureIndexes(ctx, inventoryCol); err != nil {
			panic("failed to ensure indexes: " + err.Error())
		}
		d.inventoryRepository = inventoryRepository.NewPartRepository(inventoryCol)
	}

	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongodb.Client {
	if d.mongoDBClient == nil {
		client, err := mongodb.NewClient(ctx, mongodb.Config{
			URI:             config.AppConfig().Mongo.URL(),
			Database:        config.AppConfig().Mongo.DbName(),
			ConnectTimeout:  10 * time.Second,
			MaxPoolSize:     100,
			MinPoolSize:     2,
			MaxConnIdleTime: 10 * time.Second,
		})
		if err != nil {
			panic(fmt.Sprintf("failed to connect to MongoDB: %s\n", err.Error()))
		}

		closer.AddNamed("mongodb client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}
