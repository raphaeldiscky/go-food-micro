package rabbitmq

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1/events/integrationevents"
)

func ConfigProductsRabbitMQ(
	builder configurations.RabbitMQConfigurationBuilder,
) {
	builder.AddProducer(
		integrationevents.ProductCreatedV1{},
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder) {
		},
	)
}
