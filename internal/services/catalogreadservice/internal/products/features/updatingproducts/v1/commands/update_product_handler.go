// Package commands contains the update product handler.
package commands

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
)

// UpdateProductHandler is a struct that contains the update product handler.
type UpdateProductHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

// NewUpdateProductHandler creates a new UpdateProductHandler.
func NewUpdateProductHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *UpdateProductHandler {
	return &UpdateProductHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

// Handle is a method that handles the update product command.
func (c *UpdateProductHandler) Handle(
	ctx context.Context,
	command *UpdateProduct,
) (*mediatr.Unit, error) {
	product, err := c.mongoRepository.GetProductByProductID(
		ctx,
		command.ProductID.String(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in fetching product with productID %s in the mongo repository",
				command.ProductID,
			),
		)
	}

	if product == nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"product with productID %s not found",
				command.ProductID,
			),
		)
	}

	product.Price = command.Price
	product.Name = command.Name
	product.Description = command.Description
	product.UpdatedAt = command.UpdatedAt

	_, err = c.mongoRepository.UpdateProduct(ctx, product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in updating product in the mongo repository",
		)
	}

	err = c.redisRepository.PutProduct(ctx, product.ID, product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in updating product in the redis repository",
		)
	}

	c.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} updated",
			product.ID,
		),
		logger.Fields{"ProductID": command.ProductID, "ID": product.ID},
	)

	return &mediatr.Unit{}, nil
}
