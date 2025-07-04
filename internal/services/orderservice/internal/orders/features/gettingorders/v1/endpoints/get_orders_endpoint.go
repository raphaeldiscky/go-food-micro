// Package endpoints contains the endpoints for the get orders.
package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/queries"
)

// getOrdersEndpoint is the endpoint for the get orders.
type getOrdersEndpoint struct {
	params.OrderRouteParams
}

// NewGetOrdersEndpoint creates a new get orders endpoint.
func NewGetOrdersEndpoint(p params.OrderRouteParams) route.Endpoint {
	return &getOrdersEndpoint{OrderRouteParams: p}
}

// MapEndpoint maps the endpoint.
func (ep *getOrdersEndpoint) MapEndpoint() {
	ep.OrdersGroup.GET("", ep.handler())
}

// GetAllOrders
// @Tags Orders
// @Summary Get all orders
// @Description Get all orders
// @Accept json
// @Produce json
// @Param getOrdersRequestDto query dtos.GetOrdersRequestDto false "GetOrdersRequestDto"
// @Success 200 {object} dtos.GetOrdersResponseDto
// @Router /api/v1/orders [get].
func (ep *getOrdersEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.OrdersMetrics.HTTPMetrics.GetOrdersHTTPRequests.Add(ctx, 1)

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getOrdersEndpoint_handler.GetListQueryFromCtx] error in getting data from query string",
			)
			ep.Logger.Errorf(
				fmt.Sprintf(
					"[getOrdersEndpoint_handler.GetListQueryFromCtx] err: %v",
					badRequestErr,
				),
			)

			return err
		}

		request := &dtos.GetOrdersRequestDto{ListQuery: listQuery}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getOrdersEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(fmt.Sprintf("[getOrdersEndpoint_handler.Bind] err: %v", badRequestErr))

			return badRequestErr
		}

		query := queries.NewGetOrders(request.ListQuery)

		queryResult, err := mediatr.Send[*queries.GetOrders, *dtos.GetOrdersResponseDto](ctx, query)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getOrdersEndpoint_handler.Send] error in sending GetOrders",
			)
			ep.Logger.Error(fmt.Sprintf("[getOrdersEndpoint_handler.Send] err: {%v}", err))

			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
