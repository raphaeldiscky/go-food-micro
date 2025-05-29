package queries

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/getting_orders/v1/dtos"
)

type GetOrdersHandler struct {
	log                      logger.Logger
	mongoOrderReadRepository repositories.OrderMongoRepository
	tracer                   tracing.AppTracer
}

func NewGetOrdersHandler(
	log logger.Logger,
	mongoOrderReadRepository repositories.OrderMongoRepository,
	tracer tracing.AppTracer,
) *GetOrdersHandler {
	return &GetOrdersHandler{
		log:                      log,
		mongoOrderReadRepository: mongoOrderReadRepository,
		tracer:                   tracer,
	}
}

func (c *GetOrdersHandler) Handle(
	ctx context.Context,
	query *GetOrders,
) (*dtos.GetOrdersResponseDto, error) {
	products, err := c.mongoOrderReadRepository.GetAllOrders(ctx, query.ListQuery)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[GetOrdersHandler_Handle.GetAllOrders] error in getting orders in the repository",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosV1.OrderReadDto](products)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"[GetOrdersHandler_Handle.ListResultToListResultDto] error in the mapping ListResultToListResultDto",
		)
	}

	c.log.Info("[GetOrdersHandler.Handle] orders fetched")

	return &dtos.GetOrdersResponseDto{Orders: listResultDto}, nil
}
