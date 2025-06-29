//go:build integration
// +build integration

// Package bus provides the rabbitmq bus.
package bus

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	uuid "github.com/satori/go.uuid"

	messageConsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	types3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer/json"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	rabbitmqconsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	rabbitmqproducer "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging/consumer"
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
)

var Logger = defaultlogger.GetLogger()

// TestAddRabbitMQ tests the add rabbitmq.
func TestAddRabbitMQ(t *testing.T) {
	testUtils.SkipCI(t)
	ctx := context.Background()

	// Use only one consumer to simplify the test
	fakeConsumer := consumer.NewRabbitMQFakeTestConsumerHandler[ProducerConsumerMessage]()

	serializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(Logger).
		PopulateContainerOptions(ctx, t)
	require.NoError(t, err)

	options := &config.RabbitmqOptions{
		RabbitmqHostOptions: rabbitmqHostOption,
	}

	conn, err := types.NewRabbitMQConnection(options)
	require.NoError(t, err)

	consumerFactory := rabbitmqconsumer.NewConsumerFactory(
		options,
		conn,
		serializer,
		Logger,
	)
	producerFactory := rabbitmqproducer.NewProducerFactory(
		options,
		conn,
		serializer,
		Logger,
	)

	b, err := NewRabbitmqBus(
		Logger,
		consumerFactory,
		producerFactory,
		func(builder configurations.RabbitMQConfigurationBuilder) {
			builder.AddProducer(
				ProducerConsumerMessage{},
				func(_ producerConfigurations.RabbitMQProducerConfigurationBuilder) {
				},
			)
			builder.AddConsumer(
				ProducerConsumerMessage{},
				func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
					builder.WithHandlers(
						func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
							consumerHandlerBuilder.AddHandler(fakeConsumer)
						},
					)
				},
			)
		},
	)

	require.NoError(t, err)

	err = b.Start(ctx)
	require.NoError(t, err)

	Logger.Info("Publishing message...")
	err = b.PublishMessage(
		context.Background(),
		&ProducerConsumerMessage{
			Data:    "test message data",
			Message: types3.NewMessage(uuid.NewV4().String()),
		},
		nil,
	)
	require.NoError(t, err)
	Logger.Info("Message published successfully")

	Logger.Info("Waiting for consumer to handle message...")
	err = testUtils.WaitUntilConditionMet(func() bool {
		handled := fakeConsumer.IsHandled()
		Logger.Infof("Consumer handled: %v", handled)
		return handled
	})
	assert.NoError(t, err)

	err = b.Stop()
	require.NoError(t, err)
}

// ProducerConsumerMessage is the message for the producer consumer.
type ProducerConsumerMessage struct {
	*types3.Message
	Data string
}

// NewProducerConsumerMessage creates a new producer consumer message.
func NewProducerConsumerMessage(data string) *ProducerConsumerMessage {
	return &ProducerConsumerMessage{
		Data:    data,
		Message: types3.NewMessage(uuid.NewV4().String()),
	}
}
