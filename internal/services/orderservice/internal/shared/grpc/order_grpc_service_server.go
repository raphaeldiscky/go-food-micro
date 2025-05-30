// Package grpc contains the order grpc service server.
package grpc

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	attribute2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	utils2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	uuid "github.com/satori/go.uuid"
	api "go.opentelemetry.io/otel/metric"

	dtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/dtos/v1"
	createOrderCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/commands"
	createOrderDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
	GetOrderByIDDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/dtos"
	GetOrderByIDQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/queries"
	getOrdersDtosV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/dtos"
	getOrdersQueryV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/queries"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/contracts"
	grpcOrderService "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc/genproto"
)

// OrderGrpcServiceServer is the order grpc service server.
type OrderGrpcServiceServer struct {
	ordersMetrics *contracts.OrdersMetrics
	logger        logger.Logger
	validator     *validator.Validate
}

// getGrpcMetricsAttributes returns the gRPC metrics attributes.
func getGrpcMetricsAttributes() attribute.KeyValue {
	return attribute.Key("MetricsType").String("Grpc")
}

// NewOrderGrpcService creates a new order grpc service.
func NewOrderGrpcService(
	log logger.Logger,
	val *validator.Validate,
	ordersMetrics *contracts.OrdersMetrics,
) *OrderGrpcServiceServer {
	return &OrderGrpcServiceServer{
		ordersMetrics: ordersMetrics,
		logger:        log,
		validator:     val,
	}
}

// CreateOrder creates a new order.
func (o OrderGrpcServiceServer) CreateOrder(
	ctx context.Context,
	req *grpcOrderService.CreateOrderReq,
) (*grpcOrderService.CreateOrderRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))
	o.ordersMetrics.GrpcMetrics.CreateOrderGrpcRequests.Add(
		ctx,
		1,
		api.WithAttributes(getGrpcMetricsAttributes()),
	)

	shopItemsDtos, err := mapper.Map[[]*dtosV1.ShopItemDto](req.GetShopItems())
	if err != nil {
		return nil, err
	}

	command, err := createOrderCommandV1.NewCreateOrder(
		shopItemsDtos,
		req.AccountEmail,
		req.DeliveryAddress,
		req.DeliveryTime.AsTime(),
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[OrderGrpcServiceServer_CreateOrder.StructCtx] command validation failed",
		)
		o.logger.Errorf(
			fmt.Sprintf("[OrderGrpcServiceServer_CreateOrder.StructCtx] err: %v", validationErr),
		)

		return nil, validationErr
	}

	result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *createOrderDtosV1.CreateOrderResponseDto](
		ctx,
		command,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[ProductGrpcServiceServer_CreateOrder.Send] error in sending CreateOrder",
		)
		o.logger.Errorw(
			fmt.Sprintf(
				"[ProductGrpcServiceServer_CreateOrder.Send] id: {%s}, err: %v",
				command.OrderID,
				err,
			),
			logger.Fields{"ID": command.OrderID},
		)

		return nil, err
	}

	return &grpcOrderService.CreateOrderRes{OrderID: result.OrderID.String()}, nil
}

// GetOrderByID gets the order by id.
func (o OrderGrpcServiceServer) GetOrderByID(
	ctx context.Context,
	req *grpcOrderService.GetOrderByIDReq,
) (*grpcOrderService.GetOrderByIDRes, error) {
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))
	o.ordersMetrics.GrpcMetrics.GetOrderByIDGrpcRequests.Add(
		ctx,
		1,
		api.WithAttributes(getGrpcMetricsAttributes()),
	)

	orderIDUUID, err := uuid.FromString(req.ID)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] error in converting uuid",
		)
		o.logger.Errorf(
			fmt.Sprintf(
				"[OrderGrpcServiceServer_GetOrderByID.uuid.FromString] err: %v",
				badRequestErr,
			),
		)

		return nil, badRequestErr
	}

	query, err := GetOrderByIDQueryV1.NewGetOrderByID(orderIDUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.StructCtx] query validation failed",
		)
		o.logger.Errorf(
			fmt.Sprintf("[OrderGrpcServiceServer_GetOrderByID.StructCtx] err: %v", validationErr),
		)

		return nil, validationErr
	}

	queryResult, err := mediatr.Send[*GetOrderByIDQueryV1.GetOrderByID, *GetOrderByIDDtosV1.GetOrderByIDResponseDto](
		ctx,
		query,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.Send] error in sending GetOrderByID",
		)
		o.logger.Errorw(
			fmt.Sprintf(
				"[OrderGrpcServiceServer_GetOrderByID.Send] id: {%s}, err: %v",
				query.ID,
				err,
			),
			logger.Fields{"ID": query.ID},
		)

		return nil, err
	}

	q := queryResult.Order
	order, err := mapper.Map[*grpcOrderService.OrderReadModel](q)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrderByID.Map] error in mapping order",
		)

		return nil, utils2.TraceStatusFromContext(ctx, err)
	}

	return &grpcOrderService.GetOrderByIDRes{Order: order}, nil
}

// SubmitOrder submits an order.
func (o OrderGrpcServiceServer) SubmitOrder(
	_ context.Context,
	_ *grpcOrderService.SubmitOrderReq,
) (*grpcOrderService.SubmitOrderRes, error) {
	return nil, nil
}

// UpdateShoppingCart updates the shopping cart.
func (o OrderGrpcServiceServer) UpdateShoppingCart(
	_ context.Context,
	_ *grpcOrderService.UpdateShoppingCartReq,
) (*grpcOrderService.UpdateShoppingCartRes, error) {
	return nil, nil
}

// GetOrders gets the orders.
func (o OrderGrpcServiceServer) GetOrders(
	ctx context.Context,
	req *grpcOrderService.GetOrdersReq,
) (*grpcOrderService.GetOrdersRes, error) {
	o.ordersMetrics.GrpcMetrics.GetOrdersGrpcRequests.Add(
		ctx,
		1,
		api.WithAttributes(getGrpcMetricsAttributes()),
	)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(attribute2.Object("Request", req))

	query := getOrdersQueryV1.NewGetOrders(
		&utils.ListQuery{Page: int(req.Page), Size: int(req.Size)},
	)

	queryResult, err := mediatr.Send[*getOrdersQueryV1.GetOrders, *getOrdersDtosV1.GetOrdersResponseDto](
		ctx,
		query,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrders.Send] error in sending GetOrders",
		)
		o.logger.Error(fmt.Sprintf("[OrderGrpcServiceServer_GetOrders.Send] err: {%v}", err))

		return nil, err
	}

	ordersResponse, err := mapper.Map[*grpcOrderService.GetOrdersRes](queryResult.Orders)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[OrderGrpcServiceServer_GetOrders.Map] error in mapping orders",
		)

		return nil, err
	}

	return ordersResponse, nil
}
