package part

import (
	"context"
	"testing"

	"github.com/stretchr/testify/suite"

	"github.com/m4kson/rocket-factory/inventory/internal/repository/mocks"
)

type ServiceSuite struct {
	suite.Suite

	ctx context.Context

	inventoryRepository *mocks.PartRepository

	service *service
}

func (s *ServiceSuite) SetupTest() {
	s.ctx = context.Background()

	s.inventoryRepository = mocks.NewPartRepository(s.T())

	s.service = NewPartService(
		s.inventoryRepository,
	)
}

func (s *ServiceSuite) TearDownTest() {}

func TestServiceIntegration(t *testing.T) {
	suite.Run(t, new(ServiceSuite))
}
