package app

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/m4kson/rocket-factory/order/internal/client/grpc"
	inventoryClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	paymentClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
	"github.com/m4kson/rocket-factory/order/internal/config"
	"github.com/m4kson/rocket-factory/order/internal/db/postgres"
	repository "github.com/m4kson/rocket-factory/order/internal/repository"
	orderRepository "github.com/m4kson/rocket-factory/order/internal/repository/orders"
	service "github.com/m4kson/rocket-factory/order/internal/service"
	orderService "github.com/m4kson/rocket-factory/order/internal/service/orders"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	"github.com/m4kson/rocket-factory/platform/pkg/migrator"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type diContainer struct {
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient

	postgresPool    *pgxpool.Pool
	migrationRunner *migrator.Migrator

	orderRepository repository.OrderRepository
	orderService    service.OrderService
	orderV1API      orderV1.Handler
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryClient() grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		inventoryConn, err := grpc.NewClient(
			config.AppConfig().GrpcClient.InventoryGrpcAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			//log.Error("failed to create inventory grpc connection", slog.String("err", err))
		}

		closer.AddNamed("Inventory gRPC connection", func(ctx context.Context) error {
			return inventoryConn.Close()
		})

		d.inventoryClient = inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	}

	return d.inventoryClient
}

func (d *diContainer) PaymentClient() grpcClient.PaymentClient {
	if d.paymentClient == nil {
		paymentConn, err := grpc.NewClient(
			config.AppConfig().GrpcClient.PaymentGrpcAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)

		if err != nil {
			// log
		}

		closer.AddNamed("Payment gRPC connection", func(ctx context.Context) error {
			return paymentConn.Close()
		})

		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))
	}

	return d.paymentClient
}

func (d *diContainer) PostgresPool(ctx context.Context) *pgxpool.Pool {
	if d.postgresPool == nil {
		dbPort, err := strconv.Atoi(config.AppConfig().Postgres.Port())
		if err != nil {
			//log
		}

		pool, err := postgres.NewPool(ctx, postgres.Config{
			Host:     config.AppConfig().Postgres.Host(),
			Port:     dbPort,
			User:     config.AppConfig().Postgres.User(),
			Password: config.AppConfig().Postgres.Password(),
			DBName:   config.AppConfig().Postgres.DbName(),

			MaxConns:          int32(runtime.NumCPU() * 4),
			MinConns:          2,
			MaxConnLifetime:   time.Hour,
			MaxConnIdleTime:   30 * time.Minute,
			HealthCehckPeriod: time.Minute,
		})

		if err != nil {
			// log
			panic(fmt.Sprintf("failed to connect to Postgres: %s\n", err.Error()))
		}

		closer.AddNamed("Postgres pool", func(context.Context) error {
			pool.Close()
			return nil
		})
	}

	return d.postgresPool
}

func (d *diContainer) MigrationRunner() *migrator.Migrator {
	if d.migrationRunner == nil {
		migrationsPath := config.AppConfig().Postgres.MigrationsPath()

		dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.AppConfig().Postgres.Host(),
			config.AppConfig().Postgres.Port(),
			config.AppConfig().Postgres.User(),
			config.AppConfig().Postgres.Password(),
			config.AppConfig().Postgres.DbName(),
		)
		poolCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			panic(fmt.Sprintf("failed to parse Postgres config: %s\n", err.Error()))
		}

		d.migrationRunner = migrator.NewMigrator(stdlib.OpenDB(*poolCfg.ConnConfig), migrationsPath)
	}

	return d.migrationRunner
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepository.NewRepository(d.PostgresPool(ctx))
	}

	return d.orderRepository
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewOrderService(d.OrderRepository(ctx), d.PaymentClient(), d.InventoryClient())
	}

	return d.orderService
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = ordersApi.NewAPI(d.OrderService(ctx))
	}

	return d.orderV1API
}
