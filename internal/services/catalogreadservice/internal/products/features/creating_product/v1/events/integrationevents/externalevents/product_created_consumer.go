// Package externalevents contains the product created consumer.
package externalevents

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"go.opentelemetry.io/otel/attribute"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"
)

// ProductCreatedConsumer is a struct that contains the product created consumer.
type ProductCreatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

// NewProductCreatedConsumer creates a new ProductCreatedConsumer.
func NewProductCreatedConsumer(
	log logger.Logger,
	val *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &ProductCreatedConsumer{
		logger:    log,
		validator: val,
		tracer:    tracer,
	}
}

// Handle is a method that handles the product created consumer.
func (c *ProductCreatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	ctx, span := c.tracer.Start(ctx, "productCreatedConsumer.Handle")
	defer span.End()

	product, ok := consumeContext.Message().(*ProductCreatedV1)
	if !ok {
		err := errors.New("error in casting message to ProductCreatedV1")
		c.logger.Errorw("Failed to cast message", logger.Fields{"error": err})
		span.RecordError(err)

		return err
	}

	// Validate the message
	if err := c.validator.Struct(product); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"message validation failed",
		)
		c.logger.Errorw("Message validation failed", logger.Fields{"error": validationErr})
		span.RecordError(validationErr)

		return validationErr
	}

	span.SetAttributes(
		attribute.String("productId", product.ProductID),
		attribute.String("name", product.Name),
		attribute.Float64("price", product.Price),
	)

	command, err := v1.NewCreateProduct(
		product.ProductID,
		product.Name,
		product.Description,
		product.Price,
		product.CreatedAt,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"command validation failed",
		)
		c.logger.Errorw("Command validation failed", logger.Fields{"error": validationErr})
		span.RecordError(validationErr)

		return validationErr
	}

	_, err = mediatr.Send[*v1.CreateProduct, *dtos.CreateProductResponseDto](
		ctx,
		command,
	)
	if err != nil {
		err = errors.WithMessage(
			err,
			fmt.Sprintf(
				"error in sending CreateProduct with id: {%s}",
				command.ProductID,
			),
		)
		c.logger.Errorw(
			"Failed to send CreateProduct command",
			logger.Fields{
				"error":     err,
				"productId": command.ProductID,
			},
		)
		span.RecordError(err)

		return err
	}

	c.logger.Infow(
		"Product consumer handled successfully",
		logger.Fields{
			"productId": command.ProductID,
			"traceId":   span.SpanContext().TraceID().String(),
		},
	)

	return nil
}
