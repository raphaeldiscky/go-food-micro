package endpoints

import (
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

type getProductByIDEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductByIdEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductByIDEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductByIDEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product
// @Description Get product by id
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} dtos.GetProductByIdResponseDto
// @Router /api/v1/products/{id} [get]
func (ep *getProductByIDEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.GetProductByIDRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		query, err := queries.NewGetProductById(request.Id)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"query validation failed",
			)

			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.GetProductByID, *dtos.GetProductByIdResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProductByID",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
