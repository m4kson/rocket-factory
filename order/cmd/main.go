package main

import (
	"context"
	"errors"
	"fmt"
	"log"
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
	"github.com/m4kson/rocket-factory/order/internal/migrator"
	ordersRepo "github.com/m4kson/rocket-factory/order/internal/repository/orders"
	ordersService "github.com/m4kson/rocket-factory/order/internal/service/orders"
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

	ctx := context.Background()

	inventoryConn, err := grpc.NewClient(
		config.AppConfig().GrpcClient.InventoryGrpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create inventory grpc connection: %v", err)
		return
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(
		config.AppConfig().GrpcClient.PaymentGrpcAddr(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create payment grpc connection: %v", err)
		return
	}
	defer paymentConn.Close()

	invClient := inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))

	dbPort, err := strconv.Atoi(config.AppConfig().Postgres.Port())
	if err != nil {
		log.Printf("Can't read data base port from .env: %v\n", err)
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
		log.Printf("Ошибка подключения к базе данных: %v\n", err)
		return
	}
	defer pool.Close()

	log.Printf("connected to database %s on %s:%d as user %s\n", os.Getenv("POSTGRES_DB"), os.Getenv("POSTGRES_HOST"), dbPort, os.Getenv("POSTGRES_USER"))

	migrationsPath := config.AppConfig().Postgres.MigrationsPath()
	migrationRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsPath)

	err = migrationRunner.Up()
	if err != nil {
		log.Printf("migration runner failed: %v", err)
		return
	}

	orderRepository := ordersRepo.NewRepository(pool)
	orderService := ordersService.NewOrderService(orderRepository, payClient, invClient)

	ogenServer, err := orderV1.NewServer(ordersApi.NewAPI(orderService))
	if err != nil {
		log.Fatal("can't create ogen server")
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
		log.Printf("Listening on port %s", config.AppConfig().HttpServer.Port())
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listen on port %s: %v", config.AppConfig().HttpServer.Port(), err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), config.AppConfig().HttpServer.ShutdownTimeout())
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}
	log.Println("Server gracefully stopped")
}
