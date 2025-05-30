// Package mediatr contains the mediator configurations for the orderservice.
package mediatr

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	mediatr "github.com/mehdihadeli/go-mediatr"

	repositories2 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	createOrderCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/commands"
	createOrderDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
	GetOrderByIDDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/dtos"
	GetOrderByIDQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/queries"
	getOrdersDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/dtos"
	getOrdersQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/queries"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
)

// ConfigOrdersMediator configures the orders mediator.
func ConfigOrdersMediator(
	log logger.Logger,
	mongoOrderReadRepository repositories2.OrderMongoRepository,
	orderAggregateStore store.AggregateStore[*aggregate.Order],
	tracer tracing.AppTracer,
) error {
	// https://stackoverflow.com/questions/72034479/how-to-implement-generic-interfaces
	err := mediatr.RegisterRequestHandler[*createOrderCommandV1.CreateOrder, *createOrderDtosV1.CreateOrderResponseDto](
		createOrderCommandV1.NewCreateOrderHandler(log, orderAggregateStore, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*GetOrderByIDQueryV1.GetOrderByID, *GetOrderByIDDtosV1.GetOrderByIDResponseDto](
		GetOrderByIDQueryV1.NewGetOrderByIDHandler(log, mongoOrderReadRepository, tracer),
	)
	if err != nil {
		return err
	}

	err = mediatr.RegisterRequestHandler[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
		getOrdersQueryV1.NewGetOrdersHandler(log, mongoOrderReadRepository, tracer),
	)
	if err != nil {
		return err
	}

	return nil
}
