// Package queries contains the queries for the get order by id.
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
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/dtos"
)

// GetOrderByIDHandler is the handler for the get order by id.
type GetOrderByIDHandler struct {
	log                  logger.Logger
	orderMongoRepository repositories.OrderMongoRepository
	tracer               tracing.AppTracer
}

// NewGetOrderByIDHandler creates a new get order by id handler.
func NewGetOrderByIDHandler(
	log logger.Logger,
	orderMongoRepository repositories.OrderMongoRepository,
	tracer tracing.AppTracer,
) *GetOrderByIDHandler {
	return &GetOrderByIDHandler{
		log:                  log,
		orderMongoRepository: orderMongoRepository,
		tracer:               tracer,
	}
}

// Handle handles the get order by id query.
func (q *GetOrderByIDHandler) Handle(
	ctx context.Context,
	query *GetOrderByID,
) (*dtos.GetOrderByIDResponseDto, error) {
	// get order by order-read id
	order, err := q.orderMongoRepository.GetOrderByID(ctx, query.ID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"[GetOrderByIDHandler_Handle.GetOrderByID] error in getting order with id %s in the mongo repository",
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
					"[GetOrderByIDHandler_Handle.GetOrderByID] error in getting order with orderId %s in the mongo repository",
					query.ID.String(),
				),
			)
		}
	}

	orderDto, err := mapper.Map[*dtosV1.OrderReadDto](order)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[GetOrderByIDHandler_Handle.Map] error in the mapping order",
		)
	}

	q.log.Infow(
		fmt.Sprintf("[GetOrderByIDHandler.Handle] order with id: {%s} fetched", query.ID.String()),
		logger.Fields{"ID": query.ID},
	)

	return &dtos.GetOrderByIDResponseDto{Order: orderDto}, nil
}
