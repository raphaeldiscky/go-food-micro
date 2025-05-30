// Package consumer provides the rabbitmq consumer.
package consumer

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	uuid "github.com/satori/go.uuid"

	messageConsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/pipeline"
	types3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer/json"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/bus"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	rabbitmqConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/messaging/consumer"
	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
)

var Logger = defaultlogger.GetLogger()

// TestConsumerWithFakeMessage tests the consumer with fake message.
func TestConsumerWithFakeMessage(t *testing.T) {
	testUtils.SkipCI(t)

	ctx := context.Background()

	// options := &config.RabbitmqOptions{
	//	RabbitmqHostOptions: &config.RabbitmqHostOptions{
	//		UserName: "guest",
	//		Password: "guest",
	//		HostName: "localhost",
	//		Port:     5672,
	//	},
	//}

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(Logger).
		PopulateContainerOptions(ctx, t)
	require.NoError(t, err)

	options := &config.RabbitmqOptions{
		RabbitmqHostOptions: rabbitmqHostOption,
	}

	conn, err := types.NewRabbitMQConnection(options)
	require.NoError(t, err)

	eventSerializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)
	consumerFactory := NewConsumerFactory(
		options,
		conn,
		eventSerializer,
		Logger,
	)
	producerFactory := producer.NewProducerFactory(
		options,
		conn,
		eventSerializer,
		Logger,
	)

	fakeHandler := consumer.NewRabbitMQFakeTestConsumerHandler[ProducerConsumerMessage]()

	rabbitmqBus, err := bus.NewRabbitmqBus(
		Logger,
		consumerFactory,
		producerFactory,
		func(builder rabbitmqConfigurations.RabbitMQConfigurationBuilder) {
			builder.AddConsumer(
				ProducerConsumerMessage{},
				func(consumerBuilder configurations.RabbitMQConsumerConfigurationBuilder) {
					consumerBuilder.WithHandlers(
						func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
							consumerHandlerBuilder.AddHandler(fakeHandler)
						},
					)
				},
			)
		},
	)
	if err != nil {
		t.Fatalf("Error creating rabbitmq bus: %v", err)
	}

	err = rabbitmqBus.Start(ctx)
	if err != nil {
		t.Fatalf("Error starting rabbitmq bus: %v", err)
	}
	defer func() {
		if err := rabbitmqBus.Stop(); err != nil {
			t.Fatalf("Error stopping rabbitmq bus: %v", err)
		}
	}()

	time.Sleep(time.Second * 1)

	require.NoError(t, err)

	err = rabbitmqBus.PublishMessage(
		ctx,
		NewProducerConsumerMessage("test"),
		nil,
	)
	for err != nil {
		err = rabbitmqBus.PublishMessage(
			ctx,
			NewProducerConsumerMessage("test"),
			nil,
		)
	}

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeHandler.IsHandled()
	})

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

// TestMessageHandler is the test message handler.
type TestMessageHandler struct{}

// Handle handles the message.
func (t *TestMessageHandler) Handle(
	_ context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	msg, ok := consumeContext.Message().(*ProducerConsumerMessage)
	if !ok {
		return fmt.Errorf("failed to type assert message to *ProducerConsumerMessage")
	}
	fmt.Println(msg)

	return nil
}

// NewTestMessageHandler creates a new test message handler.
func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

// TestMessageHandler2 is the test message handler 2.
type TestMessageHandler2 struct{}

// Handle handles the message.
func (t *TestMessageHandler2) Handle(
	_ context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	message := consumeContext.Message()
	fmt.Println(message)

	return nil
}

// NewTestMessageHandler2 creates a new test message handler 2.
func NewTestMessageHandler2() *TestMessageHandler2 {
	return &TestMessageHandler2{}
}

// Pipeline1 is the pipeline 1.
type Pipeline1 struct{}

// NewPipeline1 creates a new pipeline 1.
func NewPipeline1() pipeline.ConsumerPipeline {
	return &Pipeline1{}
}

// Handle handles the message.
func (p Pipeline1) Handle(
	ctx context.Context,
	consumerContext types3.MessageConsumeContext,
	next pipeline.ConsumerHandlerFunc,
) error {
	Logger.Info("PipelineBehaviourTest.Handled")
	Logger.Infof("pipeline got a message with id '%s'\n", consumerContext.Message().GeMessageId())

	err := next(ctx)
	if err != nil {
		return err
	}

	return nil
}
