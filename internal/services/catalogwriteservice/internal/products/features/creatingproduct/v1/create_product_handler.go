package v1

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mapper"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	dtosv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/models"
)

// createProductHandler is a struct that contains the create product handler.
type createProductHandler struct {
	fxparams.ProductHandlerParams
}

// NewCreateProductHandler is a constructor for the createProductHandler.
func NewCreateProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*CreateProduct, *dtos.CreateProductResponseDto] {
	return &createProductHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler is a method that registers the create product handler.
func (c *createProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*CreateProduct, *dtos.CreateProductResponseDto](
		c,
	)
}

// Handle is a method that handles the create product command.
func (c *createProductHandler) Handle(
	ctx context.Context,
	command *CreateProduct,
) (*dtos.CreateProductResponseDto, error) {
	product := &models.Product{
		ID:          command.ProductID,
		Name:        command.Name,
		Description: command.Description,
		Price:       command.Price,
		CreatedAt:   command.CreatedAt,
	}

	var createProductResult *dtos.CreateProductResponseDto

	result, err := c.ProductRepository.CreateProduct(ctx, product)
	if err != nil {
		return nil, err
	}

	productDto, err := mapper.Map[*dtosv1.ProductDto](result)
	if err != nil {
		return nil, customErrors.NewInternalServerErrorWrap(
			err,
			"error in the mapping ProductDto",
		)
	}

	productCreated := integrationevents.NewProductCreatedV1(
		productDto,
	)

	err = c.RabbitmqProducer.PublishMessage(ctx, productCreated, nil)
	if err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing ProductCreated integration_events event",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"ProductCreated message with messageId `%s` published to the rabbitmq broker",
			productCreated.MessageId,
		),
		logger.Fields{"MessageId": productCreated.MessageId},
	)

	createProductResult = &dtos.CreateProductResponseDto{
		ProductID: product.ID,
	}

	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' created",
			command.ProductID,
		),
		logger.Fields{
			"ID":        command.ProductID,
			"MessageId": productCreated.MessageId,
		},
	)

	return createProductResult, err
}
