package v1

import (
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"

	"emperror.dev/errors"
	"github.com/labstack/echo/v4"
	"github.com/mehdihadeli/go-mediatr"
)

// getProductByIDEndpoint is a struct that contains the get product by id endpoint.
type getProductByIDEndpoint struct {
	fxparams.ProductRouteParams
}

// NewGetProductByIDEndpoint is a constructor for the getProductByIDEndpoint.
func NewGetProductByIDEndpoint(
	params fxparams.ProductRouteParams,
) route.Endpoint {
	return &getProductByIDEndpoint{ProductRouteParams: params}
}

func (ep *getProductByIDEndpoint) MapEndpoint() {
	ep.ProductsGroup.GET("/:id", ep.handler())
}

// GetProductByID
// @Tags Products
// @Summary Get product by id
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

		query, err := NewGetProductByIDWithValidation(request.ProductID)
		if err != nil {
			return err
		}

		queryResult, err := mediatr.Send[*GetProductByID, *dtos.GetProductByIDResponseDto](
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
