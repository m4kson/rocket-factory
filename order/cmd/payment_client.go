package main

import (
	"context"
	"log"

	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const paymentServerAdderss = "localhost:50051"

type PaymentClient struct {
	client paymentV1.PaymentServiceClient
	conn   *grpc.ClientConn
}

func NewPaymentClient() (*PaymentClient, error) {
	conn, err := grpc.NewClient(paymentServerAdderss, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("falled to connect: %v", err)
		return nil, err
	}

	client := paymentV1.NewPaymentServiceClient(conn)

	return &PaymentClient{client: client, conn: conn}, nil
}

// Close закрывает gRPC соединение
func (pc *PaymentClient) Close() error {
	if pc.conn != nil {
		return pc.conn.Close()
	}
	return nil
}

func PayOrder(ctx context.Context, client paymentV1.PaymentServiceClient, request *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	response, err := client.PayOrder(ctx, request)
	if err != nil {
		log.Printf("failed to pay order: %v", err)
		return nil, err
	}

	return response, nil
}
