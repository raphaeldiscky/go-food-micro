package queries

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

// GetProductByIDHandler is a struct that contains the get product by id handler.
type GetProductByIDHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

// NewGetProductByIDHandler creates a new GetProductByIDHandler.
func NewGetProductByIDHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *GetProductByIDHandler {
	return &GetProductByIDHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

// getProductFromCache attempts to retrieve the product from Redis cache.
func (q *GetProductByIDHandler) getProductFromCache(
	ctx context.Context,
	id string,
) (*models.Product, error) {
	redisProduct, err := q.redisRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("error in getting product with id %s in the redis repository", id),
		)
	}

	return redisProduct, nil
}

// getProductFromMongo retrieves the product from MongoDB with fallback to ProductID.
func (q *GetProductByIDHandler) getProductFromMongo(
	ctx context.Context,
	id string,
) (*models.Product, error) {
	mongoProduct, err := q.mongoRepository.GetProductByID(ctx, id)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf("error in getting product with id %s in the mongo repository", id),
		)
	}

	if mongoProduct == nil {
		mongoProduct, err = q.mongoRepository.GetProductByProductID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return mongoProduct, nil
}

// cacheProduct stores the product in Redis cache.
func (q *GetProductByIDHandler) cacheProduct(ctx context.Context, product *models.Product) error {
	return q.redisRepository.PutProduct(ctx, product.ID, product)
}

// Handle is a method that handles the get product by id query.
func (q *GetProductByIDHandler) Handle(
	ctx context.Context,
	query *GetProductByID,
) (*dtos.GetProductByIDResponseDto, error) {
	id := query.ID.String()

	// Try to get from cache first
	product, err := q.getProductFromCache(ctx, id)
	if err != nil {
		return nil, err
	}

	// If not in cache, get from MongoDB
	if product == nil {
		product, err = q.getProductFromMongo(ctx, id)
		if err != nil {
			return nil, err
		}

		// Cache the product for future requests
		if err := q.cacheProduct(ctx, product); err != nil {
			return new(dtos.GetProductByIDResponseDto), err
		}
	}

	// Map to DTO
	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	q.log.Infow(
		fmt.Sprintf("product with id: {%s} fetched", id),
		logger.Fields{"ProductID": product.ProductID, "ID": product.ID},
	)

	return &dtos.GetProductByIDResponseDto{Product: productDto}, nil
}
