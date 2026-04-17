package payment

import (
	"context"

	"github.com/m4kson/rocket-factory/order/internal/client/converter"
	"github.com/m4kson/rocket-factory/order/internal/model"
)

func (c *client) PayOrder(ctx context.Context, requestModel model.PayOrderRequest) (model.PayOrderResponse, error) {
	request := converter.PayOrderRequestToProto(requestModel)
	response, err := c.generatedClient.PayOrder(ctx, &request)
	if err != nil {
		return model.PayOrderResponse{}, err
	}

	return converter.PayOrderResponseToModel(response), err
}
