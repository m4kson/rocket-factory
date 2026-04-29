package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"syscall"

	"github.com/m4kson/rocket-factory/payment/internal/app"
	"github.com/m4kson/rocket-factory/payment/internal/config"
	"github.com/m4kson/rocket-factory/platform/pkg/closer"
)

const configPath = "../deploy/compose/payment/.env"

func main() {
	err := config.Load(configPath)
	if err != nil {
		panic(fmt.Errorf("failed to load config: %w", err))
	}

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	a, err := app.New(context.Background())
	if err != nil {
		slog.Error("failed to initialize app", slog.String("err", err.Error()))
		os.Exit(1)
	}

	if err = a.Run(); err != nil {
		slog.Error("app exited with error", slog.String("err", err.Error()))
		os.Exit(1)
	}
}

//func main() {
//	err := config.Load(configPath)
//	if err != nil {
//		panic(fmt.Errorf("failed to load config: %w", err))
//	}
//
//	lis, err := net.Listen("tcp", ":"+config.AppConfig().Grpc.Port())
//	if err != nil {
//		log.Printf("failed to listen: %v", err)
//		return
//	}
//
//	s := grpc.NewServer()
//
//	service := paymentService.NewPaymentService()
//	api := paymentApi.NewAPI(service)
//
//	paymentV1.RegisterPaymentServiceServer(s, api)
//
//	reflection.Register(s)
//
//	go func() {
//		log.Printf("gRPC server is listening on port %s\n", config.AppConfig().Grpc.Port())
//		if err := s.Serve(lis); err != nil {
//			log.Printf("failed to serve: %v\n", err)
//			return
//		}
//	}()
//
//	quit := make(chan os.Signal, 1)
//	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
//	<-quit
//	log.Println("Shutting down gRPC server...")
//	s.GracefulStop()
//	log.Println("Server gracefully stopped")
//}
