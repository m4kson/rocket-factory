package inventory

import (
	def "github.com/m4kson/rocket-factory/order/internal/client/grpc"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

var _ def.InventoryClient = (*client)(nil)

type client struct {
	generatedClient inventoryV1.InventoryServiceClient
}

func NewClient(generatedClient inventoryV1.InventoryServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
