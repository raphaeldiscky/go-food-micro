// Package v1 contains the search products endpoint.
package v1

import (
	"net/http"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1/dtos"
)

// searchProductsEndpoint is a struct that contains the search products endpoint.
type searchProductsEndpoint struct {
	fxparams.ProductRouteParams
}

// NewSearchProductsEndpoint is a constructor for the searchProductsEndpoint.
func NewSearchProductsEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &searchProductsEndpoint{ProductRouteParams: params}
}

func (ep *searchProductsEndpoint) MapEndpoint() {
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
func (ep *searchProductsEndpoint) handler() echo.HandlerFunc {
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

		query, err := NewSearchProductsWithValidation(
			request.SearchText,
			request.ListQuery,
		)
		if err != nil {
			return err
		}

		queryResult, err := mediatr.Send[*SearchProducts, *dtos.SearchProductsResponseDto](
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
