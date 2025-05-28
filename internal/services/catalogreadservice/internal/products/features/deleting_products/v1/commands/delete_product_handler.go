// Package commands contains the delete product command handler.
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

// DeleteProductCommand is a struct that contains the delete product command.
type DeleteProductCommand struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

// NewDeleteProductHandler creates a new DeleteProductHandler.
func NewDeleteProductHandler(
	log logger.Logger,
	repository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *DeleteProductCommand {
	return &DeleteProductCommand{
		log:             log,
		mongoRepository: repository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

// Handle is a method that handles the delete product command.
func (c *DeleteProductCommand) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	product, err := c.mongoRepository.GetProductByProductID(
		ctx,
		command.ProductID.String(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in fetching product with productId %s in the mongo repository",
				command.ProductID,
			),
		)
	}
	if product == nil {
		return nil, customErrors.NewNotFoundErrorWrap(
			err,
			fmt.Sprintf(
				"product with productId %s not found",
				command.ProductID,
			),
		)
	}

	if err := c.mongoRepository.DeleteProductByID(ctx, product.ID); err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in deleting product in the mongo repository",
		)
	}

	err = c.redisRepository.DeleteProduct(ctx, product.ID)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in deleting product in the redis repository",
		)
	}

	c.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} deleted",
			product.ID,
		),
		logger.Fields{"ProductID": command.ProductID, "ID": product.ID},
	)

	return &mediatr.Unit{}, nil
}
