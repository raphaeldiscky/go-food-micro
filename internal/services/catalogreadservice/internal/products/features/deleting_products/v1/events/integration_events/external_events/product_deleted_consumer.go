package externalEvents

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"

	"emperror.dev/errors"
	"github.com/go-playground/validator"
	"github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"
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
	message, ok := consumeContext.Message().(*ProductDeletedV1)
	if !ok {
		return errors.New("error in casting message to ProductDeletedV1")
	}

	productUUID, err := uuid.FromString(message.ProductId)
	if err != nil {
		badRequestErr := customErrors.NewBadRequestErrorWrap(
			err,
			"error in the converting uuid",
		)

		return badRequestErr
	}

	command, err := commands.NewDeleteProduct(productUUID)
	if err != nil {
		validationErr := customErrors.NewValidationErrorWrap(
			err,
			"command validation failed",
		)

		return validationErr
	}

	_, err = mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](ctx, command)

	c.logger.Info("productDeletedConsumer executed successfully.")

	return err
}
