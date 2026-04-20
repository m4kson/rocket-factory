package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	ordersApi "github.com/m4kson/rocket-factory/order/internal/api/order/v1"
	inventoryClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/inventory"
	paymentClient "github.com/m4kson/rocket-factory/order/internal/client/grpc/payment"
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

	dbURI := os.Getenv("DATABASE_URL")

	pool, err := pgxpool.New(ctx, dbURI)
	if err != nil {
		log.Printf("failed to connect to postgres: %v", err)
		return
	}
	defer pool.Close()

	err = pool.Ping(ctx)
	if err != nil {
		log.Printf("База данных недоступна: %v\n", err)
		return
	}

	migrationsPath := os.Getenv("MIGRATIONS_PATH")
	migrationRunner := migrator.NewMigrator(stdlib.OpenDB(*pool.Config().ConnConfig), migrationsPath)

	err = migrationRunner.Up()
	if err != nil {
		log.Printf("migration runner failed: %v", err)
		return
	}

	orderRepository := ordersRepo.NewRepository()
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
