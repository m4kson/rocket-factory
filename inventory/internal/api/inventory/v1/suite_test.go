package v1

import (
	"context"
	"testing"

	"github.com/m4kson/rocket-factory/inventory/internal/service/mocks"
	"github.com/stretchr/testify/suite"
)

type ApiSuite struct {
	suite.Suite

	ctx context.Context

	inventoryService *mocks.PartService

	api *api
}

func (a *ApiSuite) SetupTest() {
	a.ctx = context.Background()

	a.inventoryService = mocks.NewPartService(a.T())

	a.api = NewAPI(
		a.inventoryService,
	)
}

func (a *ApiSuite) TearDownTest() {}

func TestApiIntegration(t *testing.T) {
	suite.Run(t, new(ApiSuite))
}
