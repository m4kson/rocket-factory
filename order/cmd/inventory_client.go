package main

import (
	"context"
	"log"

	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const inventoryServerAdderss = "localhost:50052"

type InventoryClient struct {
	client inventoryV1.InventoryServiceClient
	conn   *grpc.ClientConn
}

func NewInventoryClient() (*InventoryClient, error) {
	conn, err := grpc.NewClient(inventoryServerAdderss, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("falled to connect: %v", err)
		return nil, err
	}

	client := inventoryV1.NewInventoryServiceClient(conn)

	return &InventoryClient{client: client, conn: conn}, nil
}

// Close закрывает gRPC соединение
func (ic *InventoryClient) Close() error {
	if ic.conn != nil {
		return ic.conn.Close()
	}
	return nil
}

func GetPart(ctx context.Context, client inventoryV1.InventoryServiceClient, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	response, err := client.GetPart(ctx, request)
	if err != nil {
		log.Printf("failed to get part: %v", err)
		return nil, err
	}

	return response, nil
}

func ListParts(ctx context.Context, client inventoryV1.InventoryServiceClient, request *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	response, err := client.ListParts(ctx, request)
	if err != nil {
		log.Printf("failed to list parts: %v", err)
		return nil, err
	}

	return response, nil
}
