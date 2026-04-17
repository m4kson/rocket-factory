package v1

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/api/inventory/v1/helpers"
	"github.com/m4kson/rocket-factory/inventory/internal/converter"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	inventoryV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/inventory/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (a *ApiSuite) TestGetPartSuccess() {
	part := helpers.CreatePart()

	a.inventoryService.On("GetPart", a.ctx, part.PartId).Return(part, nil)

	request := &inventoryV1.GetPartRequest{
		PartUuid: part.PartId.String(),
	}

	response, err := a.api.GetPart(a.ctx, request)

	a.NoError(err)
	a.Equal(converter.PartToProto(part), response.Part)
}

func (a *ApiSuite) TestGetPartNotFound() {
	partId := uuid.New()

	request := &inventoryV1.GetPartRequest{
		PartUuid: partId.String(),
	}

	a.inventoryService.On("GetPart", a.ctx, partId).Return(model.Part{}, model.ErrPartNotFound)

	response, err := a.api.GetPart(a.ctx, request)

	a.Equal(&inventoryV1.GetPartResponse{}, response)
	a.Error(err)
	st, ok := status.FromError(err)
	a.True(ok)
	a.Equal(codes.NotFound, st.Code())
	errMes := fmt.Sprintf("part: %s not found", partId)
	a.Contains(st.Message(), errMes)
}
