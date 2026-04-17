package v1

import (
	"context"
	"errors"

	"github.com/m4kson/rocket-factory/inventory/internal/converter"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *api) GetPart(ctx context.Context, request *inventoryV1.GetPartRequest) (*inventoryV1.GetPartResponse, error) {
	partId := converter.GetPartRequestToModel(request)
	part, err := a.inventoryService.GetPart(ctx, partId)
	if err != nil {
		if errors.Is(err, model.ErrPartNotFound) {
			return &inventoryV1.GetPartResponse{}, status.Errorf(codes.NotFound, "part: %s not found", partId)
		}

		return nil, err
	}

	return &inventoryV1.GetPartResponse{
		Part: converter.PartToProto(part),
	}, nil
}
