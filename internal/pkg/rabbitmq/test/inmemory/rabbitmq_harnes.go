// Package inmemory provides the rabbitmq in memory harnesss types.
package inmemory

import (
	"context"

	consumer2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

// RabbitmqInMemoryHarnesses is a struct that contains the rabbitmq in memory harnesss.
type RabbitmqInMemoryHarnesses struct {
	publishedMessage []types.IMessage
	consumedMessage  []types.IMessage
	consumerHandlers map[types.IMessage][]consumer2.ConsumerHandler
}

// NewRabbitmqInMemoryHarnesses creates a new rabbitmq in memory harnesss.
func NewRabbitmqInMemoryHarnesses() *RabbitmqInMemoryHarnesses {
	return &RabbitmqInMemoryHarnesses{}
}

// PublishMessage publishes a message.
func (r *RabbitmqInMemoryHarnesses) PublishMessage(
	_ context.Context,
	message types.IMessage,
	_ metadata.Metadata,
) error {
	r.publishedMessage = append(r.publishedMessage, message)

	return nil
}

// PublishMessageWithTopicName publishes a message with a topic name.
func (r *RabbitmqInMemoryHarnesses) PublishMessageWithTopicName(
	_ context.Context,
	message types.IMessage,
	_ metadata.Metadata,
	_ string,
) error {
	r.publishedMessage = append(r.publishedMessage, message)

	return nil
}

// IsProduced checks if a message is produced.
func (r *RabbitmqInMemoryHarnesses) IsProduced(_ func(message types.IMessage)) {
}

// AddMessageConsumedHandler adds a message consumed handler.
func (r *RabbitmqInMemoryHarnesses) AddMessageConsumedHandler(_ func(message types.IMessage)) {
}

// Start starts the rabbitmq in memory harnesss.
func (r *RabbitmqInMemoryHarnesses) Start(_ context.Context) error {
	return nil
}

// Stop stops the rabbitmq in memory harnesss.
func (r *RabbitmqInMemoryHarnesses) Stop(_ context.Context) error {
	return nil
}

// ConnectConsumerHandler connects a consumer handler.
func (r *RabbitmqInMemoryHarnesses) ConnectConsumerHandler(
	messageType types.IMessage,
	consumerHandler consumer2.ConsumerHandler,
) error {
	r.consumerHandlers[messageType] = append(r.consumerHandlers[messageType], consumerHandler)

	return nil
}

// ConnectConsumer connects a consumer.
func (r *RabbitmqInMemoryHarnesses) ConnectConsumer(
	_ types.IMessage,
	_ consumer2.Consumer,
) error {
	return nil
}

// PublishedMessages returns the published messages.
func (r *RabbitmqInMemoryHarnesses) PublishedMessages() []types.IMessage {
	return r.publishedMessage
}

// ConsumedMessages returns the consumed messages.
func (r *RabbitmqInMemoryHarnesses) ConsumedMessages() []types.IMessage {
	return r.consumedMessage
}
