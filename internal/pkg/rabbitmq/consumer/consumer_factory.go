// Package consumer provides a set of functions for the rabbitmq consumer.
package consumer

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	serializer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/consumercontracts"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// consumerFactory is a struct that contains the consumer factory.
type consumerFactory struct {
	connection      types2.IConnection
	eventSerializer serializer.MessageSerializer
	logger          logger.Logger
	rabbitmqOptions *config.RabbitmqOptions
}

// NewConsumerFactory creates a new consumer factory.
func NewConsumerFactory(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types2.IConnection,
	eventSerializer serializer.MessageSerializer,
	l logger.Logger,
) consumercontracts.ConsumerFactory {
	return &consumerFactory{
		rabbitmqOptions: rabbitmqOptions,
		logger:          l,
		eventSerializer: eventSerializer,
		connection:      connection,
	}
}

// CreateConsumer creates a new consumer.
func (c *consumerFactory) CreateConsumer(
	consumerConfiguration *consumerConfigurations.RabbitMQConsumerConfiguration,
	isConsumedNotifications ...func(message types.IMessage),
) (consumer.Consumer, error) {
	return NewRabbitMQConsumer(
		c.rabbitmqOptions,
		c.connection,
		consumerConfiguration,
		c.eventSerializer,
		c.logger,
		isConsumedNotifications...)
}

// Connection returns the connection.
func (c *consumerFactory) Connection() types2.IConnection {
	return c.connection
}
