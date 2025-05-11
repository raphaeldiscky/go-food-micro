package consumer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
)

type ConsumerHandler interface {
	Handle(ctx context.Context, consumeContext types.MessageConsumeContext) error
}
