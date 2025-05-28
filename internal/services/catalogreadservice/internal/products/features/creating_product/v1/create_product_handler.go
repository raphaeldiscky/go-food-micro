package v1

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/models"
)

// CreateProductHandler is a struct that contains the create product handler.
type CreateProductHandler struct {
	log             logger.Logger
	mongoRepository data.ProductRepository
	redisRepository data.ProductCacheRepository
	tracer          tracing.AppTracer
}

// NewCreateProductHandler creates a new CreateProductHandler.
func NewCreateProductHandler(
	log logger.Logger,
	mongoRepository data.ProductRepository,
	redisRepository data.ProductCacheRepository,
	tracer tracing.AppTracer,
) *CreateProductHandler {
	return &CreateProductHandler{
		log:             log,
		mongoRepository: mongoRepository,
		redisRepository: redisRepository,
		tracer:          tracer,
	}
}

// Handle is a method that handles the create product command.
func (c *CreateProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	product := &models.Product{
		ID:          command.ID, // we generate id ourselves because auto generate mongo string id column with type _id is not an uuid
		ProductID:   command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	createdProduct, err := c.mongoRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in creating product in the mongo repository",
		)
	}

	err = c.redisRepository.PutProduct(ctx, createdProduct.ID, createdProduct)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in creating product in the redis repository",
		)
	}

	response := &dtos.CreateProductResponseDto{ID: createdProduct.ID}

	c.log.Infow(
		fmt.Sprintf(
			"product with id: {%s} created",
			product.ID,
		),
		logger.Fields{"ProductID": command.ProductID, "ID": product.ID},
	)

	return response, nil
}
