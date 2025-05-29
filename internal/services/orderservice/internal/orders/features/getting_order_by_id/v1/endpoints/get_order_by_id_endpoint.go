package endpoints

import (
	"fmt"
	"net/http"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/getting_order_by_id/v1/queries"
)

type getOrderByIdEndpoint struct {
	params.OrderRouteParams
}

func NewGetOrderByIdEndpoint(params params.OrderRouteParams) route.Endpoint {
	return &getOrderByIdEndpoint{OrderRouteParams: params}
}

func (ep *getOrderByIdEndpoint) MapEndpoint() {
	ep.OrdersGroup.GET("/:id", ep.handler())
}

// Get Order By ID
// @Tags Orders
// @Summary Get order by id
// @Description Get order by id
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} dtos.GetOrderByIdResponseDto
// @Router /api/v1/orders/{id} [get].
func (ep *getOrderByIdEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		ep.OrdersMetrics.GetOrderByIDHTTPRequests.Add(ctx, 1)

		request := &dtos.GetOrderByIDRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.Bind] error in the binding request",
			)
			ep.Logger.Errorf(
				fmt.Sprintf("[getProductByIdEndpoint_handler.Bind] err: %v", badRequestErr),
			)

			return badRequestErr
		}

		query, err := queries.NewGetOrderById(request.ID)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"[getProductByIdEndpoint_handler.StructCtx]  query validation failed",
			)
			ep.Logger.Errorf("[getProductByIdEndpoint_handler.StructCtx] err: %v", validationErr)

			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.GetOrderById, *dtos.GetOrderByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			err = errors.WithMessage(
				err,
				"[getProductByIdEndpoint_handler.Send] error in sending GetOrderById",
			)
			ep.Logger.Errorw(
				fmt.Sprintf(
					"[getProductByIdEndpoint_handler.Send] id: {%s}, err: %v",
					query.ID,
					err,
				),
				logger.Fields{"ID": query.ID},
			)

			return err
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
