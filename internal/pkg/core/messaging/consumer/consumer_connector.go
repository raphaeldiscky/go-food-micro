// Package consumer provides a module for the consumer.
package consumer

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// ConsumerConnector is a consumer connector.
type ConsumerConnector interface {
	// ConnectConsumerHandler Add handler to existing consumer. creates new consumer if not exist
	ConnectConsumerHandler(messageType types.IMessage, consumerHandler ConsumerHandler) error
	// ConnectConsumer Add a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type
	ConnectConsumer(messageType types.IMessage, consumer Consumer) error
}
