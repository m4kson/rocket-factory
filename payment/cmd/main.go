package main

import (
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	paymentApi "github.com/m4kson/rocket-factory/payment/internal/api/payment/v1"
	paymentService "github.com/m4kson/rocket-factory/payment/internal/service/payment"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = "50051"

func main() {
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Printf("failed to listen: %v", err)
		return
	}

	s := grpc.NewServer()

	service := paymentService.NewPaymentService()
	api := paymentApi.NewAPI(service)

	paymentV1.RegisterPaymentServiceServer(s, api)

	reflection.Register(s)

	go func() {
		log.Printf("gRPC server is listening on port %s\n", grpcPort)
		if err := s.Serve(lis); err != nil {
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
