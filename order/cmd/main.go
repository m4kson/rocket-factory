package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/stdlib"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	paymentClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
	"github.com/m4kson/rocket-factory/order/internal/config"
	"github.com/m4kson/rocket-factory/order/internal/db/postgres"
	ordersRepo "github.com/m4kson/rocket-factory/order/internal/repository/orders"
	ordersService "github.com/m4kson/rocket-factory/order/internal/service/orders"
	logger "github.com/m4kson/rocket-factory/platform/pkg/logger/slogLog"
	"github.com/m4kson/rocket-factory/platform/pkg/migrator"
	orderV1 "github.com/m4kson/rocket-factory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const configPath = "../deploy/compose/order/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	log := logger.New(logger.Config{
		Level:       config.AppConfig().Logger.Level(),
		AsJson:      config.AppConfig().Logger.AsJson(),
		ServiceName: "order",
		Environment: "local", //todo add this ot config
		AddSource:   true,    //todo getEnv("ENV", "production") == "local"
	})

	ctx := context.Background()

	log.Info("initializing service")

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().GrpcClient.InventoryGrpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to create inventory grpc client", slog.String("error", err.Error()))
		return
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(
		config.AppConfig().GrpcClient.PaymentGrpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Error("failed to create payment grpc client", slog.String("error", err.Error()))
		return
	}
	defer paymentConn.Close()

	invClient := inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))

	dbPort, err := strconv.Atoi(config.AppConfig().Postgres.Port())
	if err != nil {
		log.Error("failed to parse db port", slog.String("error", err.Error()))
		return
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
		log.Error("failed to create postgres pool", slog.String("db name", config.AppConfig().Postgres.DbName()), slog.String("error", err.Error()))
		return
	}
	defer pool.Close()

	log.Info("successfully connected to postgres", slog.String("db name", config.AppConfig().Postgres.DbName()), slog.String("port", config.AppConfig().Postgres.Port()))

	migrationsPath := config.AppConfig().Postgres.MigrationsPath()
	migrationRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsPath)

	err = migrationRunner.Up()
	if err != nil {
		log.Error("failed to run migrations", slog.String("error", err.Error()))
		return
	}

	orderRepository := ordersRepo.NewRepository(pool)
	orderService := ordersService.NewOrderService(orderRepository, payClient, invClient)

	ogenServer, err := orderV1.NewServer(ordersApi.NewAPI(orderService))
	if err != nil {
		log.Error("failed to create order v1 server", slog.String("error", err.Error()))
	}

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Mount("/api/v1", ogenServer)

	server := &http.Server{
		Addr:        net.JoinHostPort("localhost", config.AppConfig().HttpServer.Port()),
		Handler:     r,
		ReadTimeout: config.AppConfig().HttpServer.ReadHeaderTimeout(),
	}

	go func() {
		log.Info("starting http server", slog.String("port", config.AppConfig().HttpServer.Port()))
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start http server", slog.String("error", err.Error()), slog.String("port", config.AppConfig().HttpServer.Port()))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().HttpServer.ShutdownTimeout())
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Error("failed to shutdown http server", slog.String("error", err.Error()))
	}
	log.Info("Server gracefully stopped")
}
