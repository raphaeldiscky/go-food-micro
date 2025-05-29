// Package endpoints contains the get product by id endpoint.
package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getproductbyid/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getproductbyid/v1/queries"
)

// getProductByIDEndpoint is a struct that contains the get product by id endpoint.
type getProductByIDEndpoint struct {
	params.ProductRouteParams
}

// NewGetProductByIDEndpoint creates a new GetProductByIDEndpoint.
func NewGetProductByIDEndpoint(
	p params.ProductRouteParams,
) route.Endpoint {
	return &getProductByIDEndpoint{
		ProductRouteParams: p,
	}
}

// MapEndpoint maps the endpoint to the router.
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
// @Success 200 {object} dtos.GetProductByIDResponseDto
// @Router /api/v1/products/{id} [get].
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

		query, err := queries.NewGetProductByID(request.ID)
		if err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"query validation failed",
			)

			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.GetProductByID, *dtos.GetProductByIDResponseDto](
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
