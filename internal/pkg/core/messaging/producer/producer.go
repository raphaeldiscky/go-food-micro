package producer

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

type Producer interface {
	PublishMessage(ctx context.Context, message types.IMessage, meta metadata.Metadata) error
	PublishMessageWithTopicName(
		ctx context.Context,
		message types.IMessage,
		meta metadata.Metadata,
		topicOrExchangeName string,
	) error
	IsProduced(func(message types.IMessage))
}
