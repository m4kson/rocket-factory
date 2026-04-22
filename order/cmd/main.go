package main

import (
	"context"
	"errors"
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
	"github.com/joho/godotenv"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	paymentClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
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

const (
	httpPort          = "8080"
	inventoryGRPCAddr = "inventory:50051"
	paymentGRPCAddr   = "payment:50052"
	// Таймауты для HTTP-сервера
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 10 * time.Second
)

func main() {
	ctx := context.Background()

	inventoryConn, err := grpc.NewClient(
		inventoryGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create inventory grpc connection: %v", err)
		return
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(
		paymentGRPCAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Printf("failed to create payment grpc connection: %v", err)
		return
	}
	defer paymentConn.Close()

	invClient := inventoryClient.NewClient(inventoryV1.NewInventoryServiceClient(inventoryConn))
	payClient := paymentClient.NewClient(paymentV1.NewPaymentServiceClient(paymentConn))

	err = godotenv.Load("../deploy/compose/order/.env")
	if err != nil {
		log.Printf("Error loading .env file")
	}

	dbPort, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		log.Printf("Can't read data base port from .env: %v\n", err)
		return
	}

	pool, err := postgres.NewPool(ctx, postgres.Config{
		Host:     os.Getenv("POSTGRES_HOST"),
		Port:     dbPort,
		User:     os.Getenv("POSTGRES_USER"),
		Password: os.Getenv("POSTGRES_PASSWORD"),
		DBName:   os.Getenv("POSTGRES_DB"),

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

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
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
		Addr:        net.JoinHostPort("localhost", httpPort),
		Handler:     r,
		ReadTimeout: readHeaderTimeout,
	}

	go func() {
		log.Printf("Listening on port %s", httpPort)
		err = server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("Failed to listen on port %s: %v", httpPort, err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	err = server.Shutdown(ctx)
	if err != nil {
		log.Printf("Failed to shutdown server gracefully: %v", err)
	}
	log.Println("Server gracefully stopped")
}
