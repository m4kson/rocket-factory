package payment

import (
	def "github.com/m4kson/rocket-factory/order/internal/client/grpc"
	paymentV1 "github.com/m4kson/rocket-factory/shared/pkg/proto/payment/v1"
)

var _ def.PaymentClient = (*client)(nil)

type client struct {
	generatedClient paymentV1.PaymentServiceClient
}

func NewClient(generatedClient paymentV1.PaymentServiceClient) *client {
	return &client{
		generatedClient: generatedClient,
	}
}
