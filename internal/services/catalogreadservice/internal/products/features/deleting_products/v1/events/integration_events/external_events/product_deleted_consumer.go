// Package externalevents contains the product deleted consumer.
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
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	"go.opentelemetry.io/otel/attribute"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"
)

// ProductDeletedConsumer is a struct that contains the product deleted consumer.
type ProductDeletedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

// NewProductDeletedConsumer creates a new ProductDeletedConsumer.
func NewProductDeletedConsumer(
	log logger.Logger,
	val *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &ProductDeletedConsumer{
		logger:    log,
		validator: val,
		tracer:    tracer,
	}
}

// Handle is a method that handles the product deleted consumer.
func (c *ProductDeletedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	ctx, span := c.tracer.Start(ctx, "productDeletedConsumer.Handle")
	defer span.End()

	message, ok := consumeContext.Message().(*ProductDeletedV1)
	if !ok {
		err := errors.New("error in casting message to ProductDeletedV1")
		c.logger.Errorw("Failed to cast message", logger.Fields{"error": err})
		span.RecordError(err)

		return err
	}

	// Validate the message
	if err := c.validator.Struct(message); err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"message validation failed",
		)
		c.logger.Errorw("Message validation failed", logger.Fields{"error": validationErr})
		span.RecordError(validationErr)

		return validationErr
	}

	span.SetAttributes(
		attribute.String("productId", message.ProductID),
		attribute.String("message", fmt.Sprintf("%+v", message)),
	)

	productUUID, err := uuid.FromString(message.ProductID)
	if err != nil {
		c.logger.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[deleteProductConsumer_Consume.uuid.FromString] error in converting uuid",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[deleteProductConsumer_Consume.uuid.FromString] err: %v",
				utils.TraceErrStatusFromSpan(span, badRequestErr),
			),
		)
		span.RecordError(badRequestErr)

		return badRequestErr
	}

	command, err := commands.NewDeleteProduct(productUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[deleteProductConsumer_Consume.NewDeleteProduct] command validation failed",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[deleteProductConsumer_Consume.NewDeleteProduct] err: %v",
				utils.TraceErrStatusFromSpan(span, validationErr),
			),
		)
		span.RecordError(validationErr)

		return validationErr
	}

	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(
			err,
			fmt.Sprintf(
				"[deleteProductConsumer_Consume.Send] error in sending DeleteProduct with id: {%s}",
				command.ProductID,
			),
		)
		c.logger.Errorw(
			"Failed to send DeleteProduct command",
			logger.Fields{
				"error":     err,
				"productId": command.ProductID,
			},
		)
		span.RecordError(err)

		return err
	}

	c.logger.Infow(
		"Product deleted consumer handled successfully",
		logger.Fields{
			"productId": command.ProductID,
			"traceId":   span.SpanContext().TraceID().String(),
		},
	)

	return nil
}
