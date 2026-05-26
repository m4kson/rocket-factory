package kafka

import "github.com/m4kson/rocket-factory/order/internal/model"

type OrderAssembledDecoder interface {
	Decode(data []byte) (model.ShipAssembledEvent, error)
}
