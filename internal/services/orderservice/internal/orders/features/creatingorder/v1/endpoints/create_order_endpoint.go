// Package createorderendpointv1 contains the create order endpoint.
package createorderendpointv1

import (
	"fmt"
	"net/http"
	"time"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/params"
	createOrderCommandV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/commands"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/dtos"
)

// createOrderEndpoint is the create order endpoint.
type createOrderEndpoint struct {
	params.OrderRouteParams
}

// NewCreateOrderEndpoint creates a new create order endpoint.
func NewCreateOrderEndpoint(p params.OrderRouteParams) route.Endpoint {
	return &createOrderEndpoint{OrderRouteParams: p}
}

// MapEndpoint maps the create order endpoint.
func (ep *createOrderEndpoint) MapEndpoint() {
	ep.OrdersGroup.POST("", ep.handler())
}

// Create Order
// @Tags Orders
// @Summary Create order
// @Description Create new order
// @Accept json
// @Produce json
// @Param CreateOrderRequestDto body dtos.CreateOrderRequestDto true "Order data"
// @Success 201 {object} dtos.CreateOrderResponseDto
// @Router /api/v1/orders [post].
func (ep *createOrderEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.OrdersMetrics.HTTPMetrics.CreateOrderHTTPRequests.Add(ctx, 1)

		request := &dtos.CreateOrderRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[createOrderEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createOrderEndpoint_handler.Bind] err: %v", badRequestErr),
			)

			return badRequestErr
		}

		command, err := createOrderCommandV1.NewCreateOrder(
			request.ShopItems,
			request.AccountEmail,
			request.DeliveryAddress,
			time.Time(request.DeliveryTime),
		)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[createOrderEndpoint_handler.StructCtx] command validation failed",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[createOrderEndpoint_handler.StructCtx] err: %v", validationErr),
			)

			return validationErr
		}

		result, err := mediatr.Send[*createOrderCommandV1.CreateOrder, *dtos.CreateOrderResponseDto](
			ctx,
			command,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[createOrderEndpoint_handler.Send] error in sending CreateOrder",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[createOrderEndpoint_handler.Send] id: {%s}, err: %v",
					command.OrderID,
					err,
				),
				logger.Fields{"ID": command.OrderID},
			)

			return err
		}

		return c.JSON(http.StatusCreated, result)
	}
}
