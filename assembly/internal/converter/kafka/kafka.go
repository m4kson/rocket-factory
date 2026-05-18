package kafka

import "github.com/m4kson/rocket-factory/assembly/internal/model"

type OrderPaidDecoder interface {
	Decode(data []byte) (model.OrderPaidEvent, error)
}
