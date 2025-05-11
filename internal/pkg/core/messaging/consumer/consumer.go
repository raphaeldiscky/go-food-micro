package consumer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

type Consumer interface {
	Start(ctx context.Context) error
	Stop() error
	ConnectHandler(handler ConsumerHandler)
	IsConsumed(func(message types.IMessage))
	GetName() string
}
