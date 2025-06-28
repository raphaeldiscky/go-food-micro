// Package rabbitmq contains the rabbitmq configurations.
package rabbitmq

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	rabbitmqConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"

	createProductExternalEventV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/creatingproduct/v1/events/integrationevents/externalevents"
	deleteProductExternalEventV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deletingproducts/v1/events/integrationevents/externalevents"
	updateProductExternalEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/updatingproducts/v1/events/integrationevents/externalevents"
)

// ConfigProductsRabbitMQ configures the rabbitmq for the products.
func ConfigProductsRabbitMQ(
	builder rabbitmqConfigurations.RabbitMQConfigurationBuilder,
	log logger.Logger,
	val *validator.Validate,
	tracer tracing.AppTracer,
) {
	// Create message instances
	productCreatedMsg := &createProductExternalEventV1.ProductCreatedV1{}
	productDeletedMsg := &deleteProductExternalEventV1.ProductDeletedV1{}
	productUpdatedMsg := &updateProductExternalEventsV1.ProductUpdatedV1{}

	// Register message types using the standard utility function
	messageTypesMap := map[string]types.IMessage{
		productCreatedMsg.GetMessageTypeName(): productCreatedMsg,
		productDeletedMsg.GetMessageTypeName(): productDeletedMsg,
		productUpdatedMsg.GetMessageTypeName(): productUpdatedMsg,
	}

	utils.RegisterCustomMessageTypesToRegistry(messageTypesMap)

	log.Infow("Registered message types for products using standard utility", logger.Fields{
		"productCreated": productCreatedMsg.GetMessageTypeName(),
		"productDeleted": productDeletedMsg.GetMessageTypeName(),
		"productUpdated": productUpdatedMsg.GetMessageTypeName(),
	})

	builder.
		AddConsumer(
			productCreatedMsg,
			func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							createProductExternalEventV1.NewProductCreatedConsumer(
								log,
								val,
								tracer,
							),
						)
					},
				)
			}).
		AddConsumer(
			productDeletedMsg,
			func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							deleteProductExternalEventV1.NewProductDeletedConsumer(
								log,
								val,
								tracer,
							),
						)
					},
				)
			}).
		AddConsumer(
			productUpdatedMsg,
			func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
				builder.WithHandlers(
					func(handlersBuilder consumer.ConsumerHandlerConfigurationBuilder) {
						handlersBuilder.AddHandler(
							updateProductExternalEventsV1.NewProductUpdatedConsumer(
								log,
								val,
								tracer,
							),
						)
					},
				)
			})
}
