package createOrderCommandV1

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/valueobject"
)

type CreateOrderHandler struct {
	log logger.Logger
	// goland can't detect this generic type, but it is ok in vscode
	aggregateStore store.AggregateStore[*aggregate.Order]
	tracer         tracing.AppTracer
}

func NewCreateOrderHandler(
	log logger.Logger,
	aggregateStore store.AggregateStore[*aggregate.Order],
	tracer tracing.AppTracer,
) *CreateOrderHandler {
	return &CreateOrderHandler{log: log, aggregateStore: aggregateStore, tracer: tracer}
}

func (c *CreateOrderHandler) Handle(
	ctx context.Context,
	command *CreateOrder,
) (*dtos.CreateOrderResponseDto, error) {
	shopItems, err := mapper.Map[[]*valueobject.ShopItem](command.ShopItems)
	if err != nil {
		return nil,
			customErrors.NewApplicationErrorWrap(
				err,
				"[CreateOrderHandler_Handle.Map] error in the mapping shopItems",
			)
	}

	order, err := aggregate.NewOrder(
		command.OrderId,
		shopItems,
		command.AccountEmail,
		command.DeliveryAddress,
		command.DeliveryTime,
		command.CreatedAt,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[CreateOrderHandler_Handle.NewOrder] error in creating new order",
		)
	}

	_, err = c.aggregateStore.Store(order, nil, ctx)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[CreateOrderHandler_Handle.Store] error in storing order aggregate",
		)
	}

	response := &dtos.CreateOrderResponseDto{OrderId: order.ID()}

	c.log.Infow(
		fmt.Sprintf("[CreateOrderHandler.Handle] order with id: {%s} created", command.OrderId),
		logger.Fields{"ProductID": command.OrderId},
	)

	return response, nil
}
