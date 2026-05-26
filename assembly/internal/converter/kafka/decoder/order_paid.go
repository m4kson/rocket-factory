package decoder

import (
	"fmt"

	"github.com/m4kson/rocket-factory/assembly/internal/model"
	events_v1 "github.com/m4kson/rocket-factory/shared/pkg/proto/events/v1"
	"google.golang.org/protobuf/proto"
)

type decoder struct{}

func NewOrderPaidDecoder() *decoder { return &decoder{} }

func (d *decoder) Decode(data []byte) (model.OrderPaidEvent, error) {
	var pb events_v1.OrderPaid
	if err := proto.Unmarshal(data, &pb); err != nil {
		return model.OrderPaidEvent{}, fmt.Errorf("failed to unmarshal protobuf: %w", err)
	}

	return model.OrderPaidEvent{
		ID:            pb.EventUuid,
		OrderID:       pb.OrderUuid,
		UserID:        pb.UserUuid,
		PaymentMethod: pb.PaymentMethod,
		TransactionID: pb.TransactionUuid,
	}, nil
}
