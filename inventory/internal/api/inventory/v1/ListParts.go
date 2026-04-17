package v1

import (
	"context"

	"github.com/m4kson/rocket-factory/inventory/internal/converter"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (a *api) ListParts(ctx context.Context, req *inventoryV1.ListPartsRequest) (*inventoryV1.ListPartsResponse, error) {
	filters := req.GetFilter()
	parts, err := a.inventoryService.ListParts(ctx, converter.FilterToModel(filters))
	if err != nil {
		return nil, err
	}

	var protoParts []*inventoryV1.Part
	for _, part := range parts {
		protoParts = append(protoParts, converter.PartToProto(part))
	}

	return &inventoryV1.ListPartsResponse{
		Parts: protoParts,
	}, nil
}
