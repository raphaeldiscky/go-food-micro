package configurations

import (
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfiguration struct {
	ProducersConfigurations []*producerConfigurations.RabbitMQProducerConfiguration
	ConsumersConfigurations []*consumerConfigurations.RabbitMQConsumerConfiguration
}
