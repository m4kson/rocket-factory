package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	inventoryAPI "github.com/m4kson/rocket-factory/inventory/internal/api/inventory/v1"
	"github.com/m4kson/rocket-factory/inventory/internal/config"
	mongodb "github.com/m4kson/rocket-factory/inventory/internal/db/mongo"
	inventoryRepository "github.com/m4kson/rocket-factory/inventory/internal/repository/part"
	inventoryService "github.com/m4kson/rocket-factory/inventory/internal/service/part"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const configPath = "../deploy/compose/inventory/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	ctx := context.Background()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", config.AppConfig().Grpc.Port()))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}

	s := grpc.NewServer()

	mongoClient, err := mongodb.NewClient(ctx, mongodb.Config{
		URI:             config.AppConfig().Mongo.URL(),
		Database:        config.AppConfig().Mongo.DbName(),
		ConnectTimeout:  10 * time.Second,
		MaxPoolSize:     100,
		MinPoolSize:     2,
		MaxConnIdleTime: 10 * time.Second,
	})

	if err != nil {
		log.Printf("failed to connect to MongoDB: %v\n", err)
		return
	}
	defer mongoClient.Disconnect(ctx)

	log.Print("connected to mongodb")

	inventoryCol := mongoClient.Collection("inventory")
	if err = mongodb.EnsureIndexes(ctx, inventoryCol); err != nil {
		log.Printf("failed to ensure indexes: %v\n", err)
		return
	}

	repo := inventoryRepository.NewPartRepository(inventoryCol)
	service := inventoryService.NewPartService(repo)
	api := inventoryAPI.NewAPI(service)

	inventoryV1.RegisterInventoryServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("grpc server listening on %s\n", config.AppConfig().Grpc.Port())
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("Server gracefully stopped")
}
