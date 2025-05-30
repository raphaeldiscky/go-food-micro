// Package producer provides a set of functions for the rabbitmq producer.
package producer

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	serializer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/producercontracts"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// producerFactory is a struct that contains the producer factory.
type producerFactory struct {
	connection      types2.IConnection
	logger          logger.Logger
	eventSerializer serializer.MessageSerializer
	rabbitmqOptions *config.RabbitmqOptions
}

// NewProducerFactory creates a new producer factory.
func NewProducerFactory(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types2.IConnection,
	eventSerializer serializer.MessageSerializer,
	l logger.Logger,
) producercontracts.ProducerFactory {
	return &producerFactory{
		rabbitmqOptions: rabbitmqOptions,
		logger:          l,
		connection:      connection,
		eventSerializer: eventSerializer,
	}
}

// CreateProducer creates a new producer.
func (p *producerFactory) CreateProducer(
	rabbitmqProducersConfiguration map[string]*producerConfigurations.RabbitMQProducerConfiguration,
	isProducedNotifications ...func(message types.IMessage),
) (producer.Producer, error) {
	return NewRabbitMQProducer(
		p.rabbitmqOptions,
		p.connection,
		rabbitmqProducersConfiguration,
		p.logger,
		p.eventSerializer,
		isProducedNotifications...)
}
