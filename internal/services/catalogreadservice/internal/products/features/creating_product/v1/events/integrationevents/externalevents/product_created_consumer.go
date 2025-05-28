package externalEvents

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	v1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creating_product/v1/dtos"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	mediatr "github.com/mehdihadeli/go-mediatr"
	"go.opentelemetry.io/otel/attribute"
)

type productCreatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductCreatedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productCreatedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productCreatedConsumer) Handle(
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
		attribute.String("productId", product.ProductId),
		attribute.String("name", product.Name),
		attribute.Float64("price", product.Price),
	)

	command, err := v1.NewCreateProduct(
		product.ProductId,
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
				command.ProductId,
			),
		)
		c.logger.Errorw(
			"Failed to send CreateProduct command",
			logger.Fields{
				"error":     err,
				"productId": command.ProductId,
			},
		)
		span.RecordError(err)
		return err
	}

	c.logger.Infow(
		"Product consumer handled successfully",
		logger.Fields{
			"productId": command.ProductId,
			"traceId":   span.SpanContext().TraceID().String(),
		},
	)

	return nil
}
