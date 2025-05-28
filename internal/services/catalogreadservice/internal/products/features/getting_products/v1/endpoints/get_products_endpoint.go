package endpoints

import (
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/params"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getting_products/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getting_products/v1/queries"

	"emperror.dev/errors"
	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
)

type getProductsEndpoint struct {
	params.ProductRouteParams
}

func NewGetProductsEndpoint(
	params params.ProductRouteParams,
) route.Endpoint {
	return &getProductsEndpoint{
		ProductRouteParams: params,
	}
}

func (ep *getProductsEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("", ep.handler())
}

// GetAllProducts
// @Tags Products
// @Summary Get all product
// @Description Get all products
// @Accept json
// @Produce json
// @Param getProductsRequestDto query dtos.GetProductsRequestDto false "GetProductsRequestDto"
// @Success 200 {object} dtos.GetProductsResponseDto
// @Router /api/v1/products [get]
func (ep *getProductsEndpoint) handler() echo.HandlerFunc {
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

		request := queries.NewGetProducts(listQuery)
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}
		query := &queries.GetProducts{ListQuery: request.ListQuery}

		queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
			ctx,
			query,
		)
		if err != nil {
			return errors.WithMessage(
				err,
				"error in sending GetProducts",
			)
		}

		return c.JSON(http.StatusOK, queryResult)
	}
}
