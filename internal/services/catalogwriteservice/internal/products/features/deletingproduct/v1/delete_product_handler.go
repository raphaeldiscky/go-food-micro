package v1

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/dtos/v1/fxparams"
	integrationEvents "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1/events/integrationevents"
)

// deleteProductHandler is a struct that contains the delete product handler.
type deleteProductHandler struct {
	fxparams.ProductHandlerParams
}

// NewDeleteProductHandler is a constructor for the deleteProductHandler.
func NewDeleteProductHandler(
	params fxparams.ProductHandlerParams,
) cqrs.RequestHandlerWithRegisterer[*DeleteProduct, *mediatr.Unit] {
	return &deleteProductHandler{
		ProductHandlerParams: params,
	}
}

// RegisterHandler is a method that registers the delete product handler.
func (c *deleteProductHandler) RegisterHandler() error {
	return mediatr.RegisterRequestHandler[*DeleteProduct, *mediatr.Unit](
		c,
	)
}

// Handle is a method that handles the delete product command.
func (c *deleteProductHandler) Handle(
	ctx context.Context,
	command *DeleteProduct,
) (*mediatr.Unit, error) {
	err := gormdbcontext.DeleteDataModelByID[*datamodels.ProductDataModel](
		ctx,
		c.CatalogsDBContext,
		command.ProductID,
	)
	if err != nil {
		return nil, err
	}

	productDeleted := integrationEvents.NewProductDeletedV1(
		command.ProductID.String(),
	)

	if err = c.RabbitmqProducer.PublishMessage(ctx, productDeleted, nil); err != nil {
		return nil, customErrors.NewApplicationErrorWrap(
			err,
			"error in publishing 'ProductDeleted' message",
		)
	}

	c.Log.Infow(
		fmt.Sprintf(
			"ProductDeleted message with messageId '%s' published to the rabbitmq broker",
			productDeleted.MessageId,
		),
		logger.Fields{"MessageId": productDeleted.MessageId},
	)

	c.Log.Infow(
		fmt.Sprintf(
			"product with id '%s' deleted",
			command.ProductID,
		),
		logger.Fields{"ID": command.ProductID},
	)

	return &mediatr.Unit{}, err
}
