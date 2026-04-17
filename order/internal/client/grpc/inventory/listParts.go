package inventory

import (
	"context"

	"github.com/m4kson/rocket-factory/order/internal/client/converter"
	"github.com/m4kson/rocket-factory/order/internal/model"
	generatedInventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
)

func (c *client) ListParts(ctx context.Context, filter model.PartsFilter) ([]model.Part, error) {
	parts, err := c.generatedClient.ListParts(ctx, &generatedInventoryV1.ListPartsRequest{
		Filter: converter.PartsFilterToProto(filter),
	})
	if err != nil {
		return nil, err
	}

	return converter.PartsListToModel(parts.Parts), nil
}
