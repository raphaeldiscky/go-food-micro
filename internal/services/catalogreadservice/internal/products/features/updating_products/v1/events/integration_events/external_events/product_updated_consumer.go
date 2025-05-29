// Package externalevents contains the product updated consumer.
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
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updating_products/v1/commands"
)

// productUpdatedConsumer is a struct that contains the product updated consumer.
type productUpdatedConsumer struct {
	logger    logger.Logger
	validator *validator.Validate
	tracer    tracing.AppTracer
}

// NewProductUpdatedConsumer creates a new ProductUpdatedConsumer.
func NewProductUpdatedConsumer(
	log logger.Logger,
	val *validator.Validate,
	tracer tracing.AppTracer,
) consumer.ConsumerHandler {
	return &productUpdatedConsumer{
		logger:    log,
		validator: val,
		tracer:    tracer,
	}
}

// Handle is a method that handles the product updated consumer.
func (c *productUpdatedConsumer) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	message, ok := consumeContext.Message().(*ProductUpdatedV1)
	if !ok {
		return errors.New("error in casting message to ProductUpdatedV1")
	}

	ctx, span := c.tracer.Start(ctx, "productUpdatedConsumer.Handle")
	span.SetAttributes(attribute.Object("Message", consumeContext.Message()))
	defer span.End()

	productUUID, err := uuid.FromString(message.ProductID)
	if err != nil {
		c.logger.WarnMsg("uuid.FromString", err)
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"[updateProductConsumer_Consume.uuid.FromString] error in the converting uuid",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.uuid.FromString] err: %v",
				utils.TraceErrStatusFromSpan(span, badRequestErr),
			),
		)

		return err
	}

	command, err := commands.NewUpdateProduct(
		productUUID,
		message.Name,
		message.Description,
		message.Price,
	)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"[updateProductConsumer_Consume.NewValidationErrorWrap] command validation failed",
		)
		c.logger.Errorf(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.StructCtx] err: {%v}",
				utils.TraceErrStatusFromSpan(span, validationErr),
			),
		)

		return err
	}

	_, err = mediatr.Send[*commands.UpdateProduct, *mediatr.Unit](ctx, command)
	if err != nil {
		err = errors.WithMessage(
			err,
			"[updateProductConsumer_Consume.Send] error in sending UpdateProduct",
		)
		c.logger.Errorw(
			fmt.Sprintf(
				"[updateProductConsumer_Consume.Send] id: {%s}, err: {%v}",
				command.ProductID,
				utils.TraceErrStatusFromSpan(span, err),
			),
			logger.Fields{"ID": command.ProductID},
		)

		return err
	}

	return nil
}
