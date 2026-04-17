package part

import (
	"github.com/google/uuid"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/service/part/helpers"
)

func (s *ServiceSuite) TestGetSuccess() {
	part := helpers.CreatePart()

	s.inventoryRepository.On("GetPart", s.ctx, part.PartId).Return(part, nil)

	res, err := s.service.GetPart(s.ctx, part.PartId)

	s.NoError(err)
	s.Equal(part, res)
}

func (s *ServiceSuite) TestGetNotFound() {
	partId := uuid.New()

	s.inventoryRepository.On("GetPart", s.ctx, partId).Return(model.Part{}, model.ErrPartNotFound)

	_, err := s.service.GetPart(s.ctx, partId)

	s.ErrorIs(err, model.ErrPartNotFound)
}
