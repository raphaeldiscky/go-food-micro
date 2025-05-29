// Package consumer provides a rabbitmq fake consumer.
package consumer

import (
	"context"
	"fmt"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/hypothesis"
)

// RabbitMQFakeTestConsumerHandler represents a rabbitmq fake test consumer handler.
type RabbitMQFakeTestConsumerHandler[T any] struct {
	isHandled  bool
	hypothesis hypothesis.Hypothesis[T]
}

// NewRabbitMQFakeTestConsumerHandlerWithHypothesis creates a new rabbitmq fake test consumer handler with a hypothesis.
func NewRabbitMQFakeTestConsumerHandlerWithHypothesis[T any](
	hypothesis hypothesis.Hypothesis[T],
) *RabbitMQFakeTestConsumerHandler[T] {
	return &RabbitMQFakeTestConsumerHandler[T]{
		hypothesis: hypothesis,
	}
}

// NewRabbitMQFakeTestConsumerHandler creates a new rabbitmq fake test consumer handler.
func NewRabbitMQFakeTestConsumerHandler[T any]() *RabbitMQFakeTestConsumerHandler[T] {
	fmt.Println("NewRabbitMQFakeTestConsumerHandler created.")

	return &RabbitMQFakeTestConsumerHandler[T]{}
}

// Handle handles a message.
func (f *RabbitMQFakeTestConsumerHandler[T]) Handle(
	ctx context.Context,
	consumeContext types.MessageConsumeContext,
) error {
	f.isHandled = true
	if f.hypothesis != nil {
		m, ok := consumeContext.Message().(T)
		if !ok {
			f.hypothesis.Test(ctx, *new(T))
		}
		f.hypothesis.Test(ctx, m)
	}

	return nil
}

// IsHandled checks if the message is handled.
func (f *RabbitMQFakeTestConsumerHandler[T]) IsHandled() bool {
	return f.isHandled
}
