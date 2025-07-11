package repositories

import (
	"context"
	"encoding/json"
	"fmt"

	"emperror.dev/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"

	redis "github.com/redis/go-redis/v9"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

const (
	redisProductPrefixKey = "product_read_service"
)

// RedisProductRepository is a struct that contains the redis product repository.
type redisProductRepository struct {
	log         logger.Logger
	redisClient redis.UniversalClient
	tracer      tracing.AppTracer
}

// NewRedisProductRepository creates a new RedisProductRepository.
func NewRedisProductRepository(
	log logger.Logger,
	redisClient redis.UniversalClient,
	tracer tracing.AppTracer,
) data.ProductCacheRepository {
	return &redisProductRepository{
		log:         log,
		redisClient: redisClient,
		tracer:      tracer,
	}
}

// PutProduct puts a product in the redis cache.
func (r *redisProductRepository) PutProduct(
	ctx context.Context,
	key string,
	product *models.Product,
) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.PutProduct")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	productBytes, err := json.Marshal(product)
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error marshaling product",
			),
		)
	}

	if err := r.redisClient.HSetNX(ctx, r.getRedisProductPrefixKey(), key, productBytes).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in updating product with key %s",
					key,
				),
			),
		)
	}

	span.SetAttributes(attribute.Object("Product", product))

	r.log.Infow(
		fmt.Sprintf(
			"product with key '%s', prefix '%s'  updated successfully",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{
			"Product":   product,
			"ID":        product.ProductID,
			"Key":       key,
			"PrefixKey": r.getRedisProductPrefixKey(),
		},
	)

	return nil
}

// GetProductByID gets a product by id from the redis cache.
func (r *redisProductRepository) GetProductByID(
	ctx context.Context,
	key string,
) (*models.Product, error) {
	ctx, span := r.tracer.Start(ctx, "redisRepository.GetProductByID")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	productBytes, err := r.redisClient.HGet(ctx, r.getRedisProductPrefixKey(), key).
		Bytes()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in getting product with Key %s from database",
					key,
				),
			),
		)
	}

	var product models.Product
	if err := json.Unmarshal(productBytes, &product); err != nil {
		return nil, utils.TraceErrStatusFromSpan(span, err)
	}

	span.SetAttributes(attribute.Object("Product", product))

	r.log.Infow(
		fmt.Sprintf(
			"product with with key '%s', prefix '%s' loaded",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{
			"Product":   product,
			"ID":        product.ProductID,
			"Key":       key,
			"PrefixKey": r.getRedisProductPrefixKey(),
		},
	)

	return &product, nil
}

// DeleteProduct deletes a product from the redis cache.
func (r *redisProductRepository) DeleteProduct(
	ctx context.Context,
	key string,
) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.DeleteProduct")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	span.SetAttributes(attribute2.String("Key", key))
	defer span.End()

	if err := r.redisClient.HDel(ctx, r.getRedisProductPrefixKey(), key).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				fmt.Sprintf(
					"error in deleting product with key %s",
					key,
				),
			),
		)
	}

	r.log.Infow(
		fmt.Sprintf(
			"product with key %s, prefix: %s deleted successfully",
			key,
			r.getRedisProductPrefixKey(),
		),
		logger.Fields{"Key": key, "PrefixKey": r.getRedisProductPrefixKey()},
	)

	return nil
}

// DeleteAllProducts deletes all products from the redis cache.
func (r *redisProductRepository) DeleteAllProducts(ctx context.Context) error {
	ctx, span := r.tracer.Start(ctx, "redisRepository.DeleteAllProducts")
	span.SetAttributes(
		attribute2.String("PrefixKey", r.getRedisProductPrefixKey()),
	)
	defer span.End()

	if err := r.redisClient.Del(ctx, r.getRedisProductPrefixKey()).Err(); err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in deleting all products",
			),
		)
	}

	r.log.Infow(
		"all products deleted",
		logger.Fields{"PrefixKey": r.getRedisProductPrefixKey()},
	)

	return nil
}

// getRedisProductPrefixKey gets the redis product prefix key.
func (r *redisProductRepository) getRedisProductPrefixKey() string {
	return redisProductPrefixKey
}
