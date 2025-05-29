// Package v1 contains the get products handler.
package v1

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/helpers/gormextensions"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtosv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// getProductsHandler is a struct that contains the get products handler.
type getProductsHandler struct {
	fxparams.ProductHandlerParams
}

// NewGetProductsHandler is a constructor for the getProductsHandler.
func NewGetProductsHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProducts, *dtos.GetProductsResponseDto] {
	return &getProductsHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler is a method that registers the get products handler.
func (c *getProductsHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProducts, *dtos.GetProductsResponseDto](
		c,
	)
}

// Handle is a method that handles the get products query.
func (c *getProductsHandler) Handle(
	ctx context.Context,
	query *GetProducts,
) (*dtos.GetProductsResponseDto, error) {
	products, err := gormextensions.Paginate[*datamodel.ProductDataModel, *models.Product](
		ctx,
		query.ListQuery,
		c.CatalogsDBContext.DB(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the fetching products",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dtosv1.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping",
		)
	}

	c.Log.Info("products fetched")

	return &dtos.GetProductsResponseDto{Products: listResultDto}, nil
}
