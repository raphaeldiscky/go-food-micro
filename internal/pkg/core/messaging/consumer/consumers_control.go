// Package consumer provides the consumers control.
package consumer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// BusControl is a interface that represents the bus control.
type BusControl interface {
	// Start starts all consumers
	Start(ctx context.Context) error
	// Stop stops all consumers
	Stop() error
	// IsConsumed checks if a message has been consumed
	IsConsumed(func(message types.IMessage))
}
