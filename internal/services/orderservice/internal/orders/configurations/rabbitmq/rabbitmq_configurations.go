// Package rabbitmq contains the rabbitmq configurations for the orderservice.
package rabbitmq

import (
	rabbitmqConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"

	createOrderIntegrationEventsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/events/integrationevents"
)

// ConfigOrdersRabbitMQ configures the orders rabbitmq.
func ConfigOrdersRabbitMQ(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
	builder.AddProducer(
		createOrderIntegrationEventsV1.OrderCreatedV1{},
		func(_ producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		})
}
