// Package bus provides the rabbitmq bus.
package bus

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	uuid "github.com/satori/go.uuid"

	messageConsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	pipeline2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/pipeline"
	types3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer/json"
	defaultlogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
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

// TestAddRabbitMQ tests the add rabbitmq.
func TestAddRabbitMQ(t *testing.T) {
	testUtils.SkipCI(t)
	ctx := context.Background()

	fakeConsumer2 := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()
	fakeConsumer3 := consumer.NewRabbitMQFakeTestConsumerHandler[*ProducerConsumerMessage]()

	serializer := json.NewDefaultMessageJsonSerializer(
		json.NewDefaultJsonSerializer(),
	)

	// rabbitmqOptions := &config.RabbitmqOptions{
	//	RabbitmqHostOptions: &config.RabbitmqHostOptions{
	//		UserName: "guest",
	//		Password: "guest",
	//		HostName: "localhost",
	//		Port:     5672,
	//	},
	//}

	rabbitmqHostOption, err := rabbitmq.NewRabbitMQTestContainers(defaultlogger.GetLogger()).
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
		defaultlogger.GetLogger(),
	)
	producerFactory := rabbitmqproducer.NewProducerFactory(
		options,
		conn,
		serializer,
		defaultlogger.GetLogger(),
	)

	b, err := NewRabbitmqBus(
		defaultlogger.GetLogger(),
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
					builder.WithHandlers(func(consumerHandlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
						consumerHandlerBuilder.AddHandler(
							NewTestMessageHandler(),
						)
						consumerHandlerBuilder.AddHandler(
							NewTestMessageHandler2(),
						)
					}).
						WIthPipelines(func(consumerPipelineBuilder pipeline2.ConsumerPipelineConfigurationBuilder) {
							consumerPipelineBuilder.AddPipeline(NewPipeline1())
						})
				},
			)
		},
	)

	require.NoError(t, err)

	// err = b.ConnectRabbitMQConsumer(ProducerConsumerMessage{}, func(consumerBuilder consumerConfigurations.RabbitMQConsumerConfigurationBuilder) {
	//	consumerBuilder.WithHandlers(func(handlerBuilder messageConsumer.ConsumerHandlerConfigurationBuilder) {
	//		handlerBuilder.AddHandler(fakeConsumer)
	//	})
	// })
	// require.NoError(t, err)

	err = b.ConnectConsumerHandler(&ProducerConsumerMessage{}, fakeConsumer2)
	require.NoError(t, err)

	err = b.ConnectConsumerHandler(&ProducerConsumerMessage{}, fakeConsumer3)
	require.NoError(t, err)

	err = b.Start(ctx)
	require.NoError(t, err)

	err = b.PublishMessage(
		context.Background(),
		&ProducerConsumerMessage{
			Data:    "ssssssssss",
			Message: types3.NewMessage(uuid.NewV4().String()),
		},
		nil,
	)
	require.NoError(t, err)

	err = testUtils.WaitUntilConditionMet(func() bool {
		return fakeConsumer2.IsHandled() && fakeConsumer3.IsHandled()
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

// TestMessageHandler is the test message handler.
type TestMessageHandler struct{}

// NewTestMessageHandler creates a new test message handler.
func NewTestMessageHandler() *TestMessageHandler {
	return &TestMessageHandler{}
}

// Handle handles the message.
func (t *TestMessageHandler) Handle(
	ctx context.Context,
	consumeContext types3.MessageConsumeContext,
) error {
	message, ok := consumeContext.Message().(*ProducerConsumerMessage)
	if !ok {
		return fmt.Errorf("failed to type assert message to *ProducerConsumerMessage")
	}
	fmt.Println(message)

	return nil
}

// TestMessageHandler2 is the test message handler 2.
type TestMessageHandler2 struct{}

// Handle handles the message.
func (t *TestMessageHandler2) Handle(
	ctx context.Context,
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
func NewPipeline1() pipeline2.ConsumerPipeline {
	return &Pipeline1{}
}

// Handle handles the message.
func (p *Pipeline1) Handle(
	ctx context.Context,
	consumerContext types3.MessageConsumeContext,
	next pipeline2.ConsumerHandlerFunc,
) error {
	fmt.Println("PipelineBehaviourTest.Handled")

	fmt.Printf(
		"pipeline got a message with id '%s'",
		consumerContext.Message().GeMessageId(),
	)

	err := next(ctx)
	if err != nil {
		return err
	}

	return nil
}
