package part

import (
	"github.com/brianvoe/gofakeit/v7"
	"github.com/m4kson/rocket-factory/inventory/internal/model"
	"github.com/m4kson/rocket-factory/inventory/internal/service/part/helpers"
)

func (s *ServiceSuite) TestListSuccess() {
	parts := helpers.CreateParts(3)

	s.inventoryRepository.On("ListParts", s.ctx, model.PartsFilter{}).Return(parts, nil)

	res, err := s.service.ListParts(s.ctx, model.PartsFilter{})

	s.NoError(err)
	s.Equal(parts, res)
}

func (s *ServiceSuite) TestListError() {
	s.inventoryRepository.On("ListParts", s.ctx, model.PartsFilter{}).Return([]model.Part{}, gofakeit.Error())

	_, err := s.service.ListParts(s.ctx, model.PartsFilter{})

	s.Error(err)
}
