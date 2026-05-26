package app

import (
	"context"
	"fmt"
	"runtime"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	grpcClient "github.com/m4kson/rocket-factory/order/internal/client/grpc"
	inventoryClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	paymentClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
	"github.com/m4kson/rocket-factory/order/internal/config"
	kafkaConverter "github.com/m4kson/rocket-factory/order/internal/converter/kafka"
	"github.com/m4kson/rocket-factory/order/internal/converter/kafka/decoder"
	"github.com/m4kson/rocket-factory/order/internal/db/postgres"
	repository "github.com/m4kson/rocket-factory/order/internal/repository"
	orderRepository "github.com/m4kson/rocket-factory/order/internal/repository/orders"
	service "github.com/m4kson/rocket-factory/order/internal/service"
	orderConsumer "github.com/m4kson/rocket-factory/order/internal/service/consumer/order_consumer"
	orderService "github.com/m4kson/rocket-factory/order/internal/service/orders"
	orderProducer "github.com/m4kson/rocket-factory/order/internal/service/producer/order_producer"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
	wrappedKafka "github.com/m4kson/rocket-factory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/m4kson/rocket-factory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/m4kson/rocket-factory/platform/pkg/kafka/producer"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	kafkaMiddleware "github.com/m4kson/rocket-factory/platform/pkg/middleware/kafka"
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

	orderProducerService service.OrderProducerService
	orderConsumerService service.OrderConsumerService

	consumerGroup sarama.ConsumerGroup
	syncProducer  sarama.SyncProducer

	orderPaidProducer      wrappedKafka.Producer
	orderAssembledConsumer wrappedKafka.Consumer
	orderAssembledDecoder  kafkaConverter.OrderAssembledDecoder
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryClient(_ context.Context) (grpcClient.InventoryClient, error) {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().GrpcClient.InventoryGrpcAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("di: inventory grpc client: %w", err)
		}

		closer.AddNamed("inventory gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.inventoryClient = inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(conn))
	}

	return d.inventoryClient, nil
}

func (d *diContainer) PaymentClient(_ context.Context) (grpcClient.PaymentClient, error) {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(
			config.AppConfig().GrpcClient.PaymentGrpcAddr(),
			grpc.WithTransportCredentials(insecure.NewCredentials()),
		)
		if err != nil {
			return nil, fmt.Errorf("di: payment grpc client: %w", err)
		}

		closer.AddNamed("payment gRPC connection", func(ctx context.Context) error {
			return conn.Close()
		})

		d.paymentClient = paymentClient.NewClient(paymentV1.NewPaymentServiceClient(conn))
	}

	return d.paymentClient, nil
}

func (d *diContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.postgresPool == nil {
		dbPort, err := strconv.Atoi(config.AppConfig().Postgres.Port())
		if err != nil {
			return nil, fmt.Errorf("di: invalid postgres port: %w", err)
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
			return nil, fmt.Errorf("di: postgres pool: %w", err)
		}

		closer.AddNamed("postgres pool", func(context.Context) error {
			pool.Close()
			return nil
		})

		d.postgresPool = pool
	}

	return d.postgresPool, nil
}

func (d *diContainer) MigrationRunner(_ context.Context) (*migrator.Migrator, error) {
	if d.migrationRunner == nil {
		dbPort, err := strconv.Atoi(config.AppConfig().Postgres.Port())
		if err != nil {
			return nil, fmt.Errorf("di: migration runner: invalid port: %w", err)
		}

		dsn := fmt.Sprintf(
			"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			config.AppConfig().Postgres.Host(),
			dbPort,
			config.AppConfig().Postgres.User(),
			config.AppConfig().Postgres.Password(),
			config.AppConfig().Postgres.DbName(),
		)

		poolCfg, err := pgxpool.ParseConfig(dsn)
		if err != nil {
			return nil, fmt.Errorf("di: migration runner: parse config: %w", err)
		}

		d.migrationRunner = migrator.NewMigrator(
			stdlib.OpenDB(*poolCfg.ConnConfig),
			config.AppConfig().Postgres.MigrationsPath(),
		)
	}

	return d.migrationRunner, nil
}

func (d *diContainer) OrderRepository(ctx context.Context) (repository.OrderRepository, error) {
	if d.orderRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order repository: %w", err)
		}
		d.orderRepository = orderRepository.NewRepository(pool)
	}

	return d.orderRepository, nil
}

func (d *diContainer) OrderService(ctx context.Context) (service.OrderService, error) {
	if d.orderService == nil {
		repo, err := d.OrderRepository(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order service: %w", err)
		}

		pay, err := d.PaymentClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order service: %w", err)
		}

		inv, err := d.InventoryClient(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order service: %w", err)
		}

		d.orderService = orderService.NewOrderService(repo, pay, inv, d.OrderProducerService(ctx))
	}

	return d.orderService, nil
}

func (d *diContainer) OrderV1API(ctx context.Context) (orderV1.Handler, error) {
	if d.orderV1API == nil {
		svc, err := d.OrderService(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order api: %w", err)
		}
		d.orderV1API = ordersApi.NewAPI(svc)
	}

	return d.orderV1API, nil
}

func (d *diContainer) OrderConsumerService(ctx context.Context) (service.OrderConsumerService, error) {
	if d.orderConsumerService == nil {
		orderSvc, err := d.OrderService(ctx)
		if err != nil {
			return nil, fmt.Errorf("di: order consumer service: %w", err)
		}

		d.orderConsumerService = orderConsumer.NewService(
			d.OrderAssembledConsumer(ctx),
			d.OrderAssembledDecoder(),
			orderSvc,
		)
	}

	return d.orderConsumerService, nil
}

func (d *diContainer) OrderProducerService(ctx context.Context) service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(
			d.OrderPaidProducer(ctx),
		)
	}
	return d.orderProducerService
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		group, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderAssembledConsumer.GroupID(),
			config.AppConfig().OrderAssembledConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return group.Close()
		})

		d.consumerGroup = group
	}

	return d.consumerGroup
}

func (d *diContainer) OrderAssembledConsumer(ctx context.Context) wrappedKafka.Consumer {
	if d.orderAssembledConsumer == nil {
		d.orderAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string{
				config.AppConfig().OrderAssembledConsumer.Topic(),
			},
			logger.FromContext(ctx),
			kafkaMiddleware.Logging(logger.FromContext(ctx)),
		)
	}

	return d.orderAssembledConsumer
}

func (d *diContainer) OrderAssembledDecoder() kafkaConverter.OrderAssembledDecoder {
	if d.orderAssembledDecoder == nil {
		d.orderAssembledDecoder = decoder.NewOrderAssembledDecoder()
	}

	return d.orderAssembledDecoder
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderPaidProducer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error {
			return p.Close()
		})

		d.syncProducer = p
	}

	return d.syncProducer
}

func (d *diContainer) OrderPaidProducer(ctx context.Context) wrappedKafka.Producer {
	if d.orderPaidProducer == nil {
		d.orderPaidProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderPaidProducer.TopicName(),
			logger.FromContext(ctx),
		)
	}
	return d.orderPaidProducer
}
