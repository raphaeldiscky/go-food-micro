package v1

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	dtoV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// GetProductByIDHandler is a struct that contains the get product by id handler.
type GetProductByIDHandler struct {
	fxparams.ProductHandlerParams
}

// NewGetProductByIDHandler is a constructor for the GetProductByIDHandler.
func NewGetProductByIDHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*GetProductByID, *dtos.GetProductByIDResponseDto] {
	return &GetProductByIDHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler is a method that registers the get product by id handler.
func (c *GetProductByIDHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*GetProductByID, *dtos.GetProductByIDResponseDto](
		c,
	)
}

// Handle is a method that handles the get product by id query.
func (c *GetProductByIDHandler) Handle(
	ctx context.Context,
	query *GetProductByID,
) (*dtos.GetProductByIDResponseDto, error) {
	product, err := gormdbcontext.FindModelByID[*datamodels.ProductDataModel, *models.Product](
		ctx,
		c.CatalogsDBContext,
		query.ProductID,
	)
	if err != nil {
		return nil, err
	}

	productDto, err := mapper.Map[*dtoV1.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.ProductID,
		),
		logger.Fields{"Id": query.ProductID.String()},
	)

	return &dtos.GetProductByIDResponseDto{Product: productDto}, nil
}
