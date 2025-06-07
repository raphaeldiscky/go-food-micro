package v1

import (
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"

	echo "github.com/labstack/echo/v4"
	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/dtos"
)

// deleteProductEndpoint is a struct that contains the delete product endpoint.
type deleteProductEndpoint struct {
	fxparams.ProductRouteParams
}

// NewDeleteProductEndpoint is a constructor for the deleteProductEndpoint.
func NewDeleteProductEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &deleteProductEndpoint{ProductRouteParams: params}
}

// MapEndpoint is a method that maps the endpoint.
func (ep *deleteProductEndpoint) MapEndpoint() {
	ep.ProductsGroup.DELETE("/:id", ep.handler())
}

// DeleteProduct
// @Tags Products
// @Summary Delete product
// @Description Delete existing product
// @Accept json
// @Produce json
// @Success 204
// @Param id path string true "Product ID"
// @Router /api/v1/products/{id} [delete].
func (ep *deleteProductEndpoint) handler() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		request := &dtos.DeleteProductRequestDto{}
		if err := c.Bind(request); err != nil {
			badRequestErr := customErrors.NewBadRequestErrorWrap(
				err,
				"error in the binding request",
			)

			return badRequestErr
		}

		command, err := NewDeleteProductWithValidation(request.ProductID)
		if err != nil {
			return err
		}

		_, err = mediatr.Send[*DeleteProduct, *mediatr.Unit](
			ctx,
			command,
		)
		if err != nil {
			return err
		}

		return c.NoContent(http.StatusNoContent)
	}
}
