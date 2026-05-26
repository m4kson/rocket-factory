package decoder

import (
	"fmt"

	"github.com/m4kson/rocket-factory/order/internal/model"
	events_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderAssembledDecoder() *decoder {
	return &decoder{}
}

func (d *decoder) Decode(data []byte) (model.ShipAssembledEvent, error) {
	var pb events_v1.ShipAssembled
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.ShipAssembledEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.ShipAssembledEvent{
		ID:           pb.EventUuid,
		OrderID:      pb.OrderUuid,
		UserID:       pb.UserUuid,
		BuildTimeSec: pb.BuildTimeSec,
	}, nil
}
