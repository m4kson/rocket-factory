package v1

import (
	"github.com/m4kson/rocket-factory/inventory/internal/service"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

type api struct {
	inventoryV1.UnimplementedInventoryServiceServer

	inventoryService service.PartService
}

func NewAPI(inventoryService service.PartService) *api {
	return &api{
		inventoryService: inventoryService,
	}
}
