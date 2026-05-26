package service

import (
	"context"

	"github.com/m4kson/rocket-factory/assembly/internal/model"
)

type AssemblyService interface{}

type ConsumerService interface {
	RunConsumer(ctx context.Context) error
}

type ShipAssembledProducerService interface {
	ProduceOrderShipAssembled(ctx context.Context, event model.ShipAssembledEvent) error
}
