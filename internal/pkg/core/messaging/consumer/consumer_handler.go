// Package consumer provides the consumer handler.
package consumer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

// ConsumerHandler is a interface that represents the consumer handler.
type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
