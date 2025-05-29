// Package consumer provides the consumer.
package consumer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// Consumer is a interface that represents the consumer.
type Consumer interface {
	// Start starts the consumer
	Start(ctx context.Context) error
	// Stop stops the consumer
	Stop() error
	// ConnectHandler connects a handler to the consumer
	ConnectHandler(handler ConsumerHandler)
	// IsConsumed checks if a message has been consumed
	IsConsumed(func(message types.IMessage))
	// GetName returns the name of the consumer
	GetName() string
}
