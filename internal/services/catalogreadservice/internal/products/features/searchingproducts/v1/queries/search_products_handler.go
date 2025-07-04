// Package queries contains the search products handler.
package queries

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/dtos"
)

// SearchProductsHandler is a struct that contains the search products handler.
type SearchProductsHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	tracer          tracing.AppTracer
}

// NewSearchProductsHandler creates a new SearchProductsHandler.
func NewSearchProductsHandler(
	log logger.Logger,
	repository data.ProductRepository,
	tracer tracing.AppTracer,
) *SearchProductsHandler {
	return &SearchProductsHandler{
		log:             log,
		mongoRepository: repository,
		tracer:          tracer,
	}
}

// Handle is a method that handles the search products query.
func (c *SearchProductsHandler) Handle(
	ctx context.Context,
	query *SearchProducts,
) (*dtos.SearchProductsResponseDto, error) {
	products, err := c.mongoRepository.SearchProducts(
		ctx,
		query.SearchText,
		query.ListQuery,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in searching products in the repository",
		)
	}

	listResultDto, err := utils.ListResultToListResultDto[*dto.ProductDto](
		products,
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping ListResultToListResultDto",
		)
	}
	c.log.Info("products fetched")

	return &dtos.SearchProductsResponseDto{Products: listResultDto}, nil
}
