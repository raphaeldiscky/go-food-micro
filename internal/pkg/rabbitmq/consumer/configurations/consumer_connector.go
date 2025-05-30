// Package configurations provides a set of functions for the rabbitmq consumer configurations.
package configurations

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// RabbitMQConsumerConnector is a interface that contains the rabbitmq consumer connector.
type RabbitMQConsumerConnector interface {
	consumer.ConsumerConnector
	// ConnectRabbitMQConsumer Add a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type
	ConnectRabbitMQConsumer(
		messageType types.IMessage,
		consumerBuilderFunc RabbitMQConsumerConfigurationBuilderFuc,
	) error
}
