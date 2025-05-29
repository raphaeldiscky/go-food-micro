package queries

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
)

type GetOrderByIdHandler struct {
	log                  logger.Logger
	orderMongoRepository repositories.OrderMongoRepository
	tracer               tracing.AppTracer
}

func NewGetOrderByIdHandler(
	log logger.Logger,
	orderMongoRepository repositories.OrderMongoRepository,
	tracer tracing.AppTracer,
) *GetOrderByIdHandler {
	return &GetOrderByIdHandler{
		log:                  log,
		orderMongoRepository: orderMongoRepository,
		tracer:               tracer,
	}
}

func (q *GetOrderByIdHandler) Handle(
	ctx context.Context,
	query *GetOrderById,
) (*dtos.GetOrderByIdResponseDto, error) {
	// get order by order-read id
	order, err := q.orderMongoRepository.GetOrderById(ctx, query.ID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"[GetOrderByIdHandler_Handle.GetProductByID] error in getting order with id %s in the mongo repository",
				query.ID.String(),
			),
		)
	}

	if order == nil {
		// get order by order-write id
		order, err = q.orderMongoRepository.GetOrderByOrderId(ctx, query.ID)
		if err != nil {
			return nil, customErrors.NewApplicationErrorWrap(
				err,
				fmt.Sprintf(
					"[GetOrderByIdHandler_Handle.GetProductByID] error in getting order with orderId %s in the mongo repository",
					query.ID.String(),
				),
			)
		}
	}

	orderDto, err := mapper.Map[*dtosV1.OrderReadDto](order)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[GetOrderByIdHandler_Handle.Map] error in the mapping order",
		)
	}

	q.log.Infow(
		fmt.Sprintf("[GetOrderByIdHandler.Handle] order with id: {%s} fetched", query.ID.String()),
		logger.Fields{"ID": query.ID},
	)

	return &dtos.GetOrderByIdResponseDto{Order: orderDto}, nil
}
