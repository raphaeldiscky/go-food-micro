package queries

import (
	"context"
	"fmt"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/dto"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

type GetProductByIdHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

func NewGetProductByIDHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *GetProductByIdHandler {
	return &GetProductByIdHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

func (q *GetProductByIdHandler) Handle(
	ctx context.Context,
	query *GetProductByID,
) (*dtos.GetProductByIDResponseDto, error) {
	redisProduct, err := q.redisRepository.GetProductByID(
		ctx,
		query.Id.String(),
	)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			fmt.Sprintf(
				"error in getting product with id %d in the redis repository",
				query.Id,
			),
		)
	}

	var product *models.Product

	if redisProduct != nil {
		product = redisProduct
	} else {
		var mongoProduct *models.Product
		mongoProduct, err = q.mongoRepository.GetProductByID(ctx, query.Id.String())
		if err != nil {
			return nil, customErrors.NewApplicationErrorWrap(err, fmt.Sprintf("error in getting product with id %d in the mongo repository", query.Id))
		}
		if mongoProduct == nil {
			mongoProduct, err = q.mongoRepository.GetProductByProductId(ctx, query.Id.String())
		}
		if err != nil {
			return nil, err
		}

		product = mongoProduct
		err = q.redisRepository.PutProduct(ctx, product.Id, product)
		if err != nil {
			return new(dtos.GetProductByIDResponseDto), err
		}
	}

	productDto, err := mapper.Map[*dto.ProductDto](product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in the mapping product",
		)
	}

	q.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} fetched",
			query.Id,
		),
		logger.Fields{"ProductId": product.ProductId, "Id": product.Id},
	)

	return &dtos.GetProductByIDResponseDto{Product: productDto}, nil
}
