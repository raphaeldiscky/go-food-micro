// Package endpoints contains the search products endpoint.
package endpoints

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/queries"
)

// SearchProductsEndpoint is a struct that contains the search products endpoint.
type SearchProductsEndpoint struct {
	params.ProductRouteParams
}

// NewSearchProductsEndpoint creates a new SearchProductsEndpoint.
func NewSearchProductsEndpoint(
	p params.ProductRouteParams,
) route.Endpoint {
	return &SearchProductsEndpoint{
		ProductRouteParams: p,
	}
}

// MapEndpoint maps the endpoint to the router.
func (ep *SearchProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/search", ep.handler())
}

// SearchProducts
// @Tags Products
// @Summary Search products
// @Description Search products
// @Accept json
// @Produce json
// @Param searchProductsRequestDto query dtos.SearchProductsRequestDto false "SearchProductsRequestDto"
// @Success 200 {object} dtos.SearchProductsResponseDto
// @Router /api/v1/products/search [get].
func (ep *SearchProductsEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		listQuery, err := utils.GetListQueryFromCtx(c)
		if err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in getting data from query string",
			)

			return badRequestErr
		}

		request := &dtos.SearchProductsRequestDto{ListQuery: listQuery}

		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		query := &queries.SearchProducts{
			SearchText: request.SearchText,
			ListQuery:  request.ListQuery,
		}

		if err := query.Validate(); err != nil {
			validationErr := customErrors.NewValidationErrorWrap(
				err,
				"query validation failed",
			)

			return validationErr
		}

		queryResult, err := mediatr.Send[*queries.SearchProducts, *dtos.SearchProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending SearchProducts",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
