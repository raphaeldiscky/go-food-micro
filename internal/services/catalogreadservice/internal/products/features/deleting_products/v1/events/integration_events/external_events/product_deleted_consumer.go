package externalEvents

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
	"go.opentelemetry.io/otel/attribute"
)

type productDeletedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

func NewProductDeletedConsumer(
	logger logger.Logger,
	validator *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productDeletedConsumer{
		logger:    logger,
		validator: validator,
		tracer:    tracer,
	}
}

func (c *productDeletedConsumer) Handle(
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
		attribute.String("productId", message.ProductId),
		attribute.String("message", fmt.Sprintf("%+v", message)),
	)

	productUUID, err := uuid.FromString(message.ProductId)
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
				command.ProductId,
			),
		)
		c.logger.Errorw(
			"Failed to send DeleteProduct command",
			logger.Fields{
				"error":     err,
				"productId": command.ProductId,
			},
		)
		span.RecordError(err)
		return err
	}

	c.logger.Infow(
		"Product deleted consumer handled successfully",
		logger.Fields{
			"productId": command.ProductId,
			"traceId":   span.SpanContext().TraceID().String(),
		},
	)

	return nil
}
