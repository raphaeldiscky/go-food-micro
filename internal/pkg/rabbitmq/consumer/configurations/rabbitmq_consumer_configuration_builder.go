// Package configurations provides a set of functions for the rabbitmq consumer configurations.
package configurations

import (
	messageConsumer "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/pipeline"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// RabbitMQConsumerConfigurationBuilder is a interface that contains the rabbitmq consumer configuration builder.
type RabbitMQConsumerConfigurationBuilder interface {
	WithHandlers(
		consumerBuilderFunc messageConsumer.ConsumerHandlerConfigurationBuilderFunc,
	) RabbitMQConsumerConfigurationBuilder
	WIthPipelines(
		pipelineBuilderFunc pipeline.ConsumerPipelineConfigurationBuilderFunc,
	) RabbitMQConsumerConfigurationBuilder
	WithExitOnError(exitOnError bool) RabbitMQConsumerConfigurationBuilder
	WithAutoAck(ack bool) RabbitMQConsumerConfigurationBuilder
	WithNoLocal(noLocal bool) RabbitMQConsumerConfigurationBuilder
	WithNoWait(noWait bool) RabbitMQConsumerConfigurationBuilder
	WithConcurrencyLimit(limit int) RabbitMQConsumerConfigurationBuilder
	WithPrefetchCount(count int) RabbitMQConsumerConfigurationBuilder
	WithConsumerID(consumerId string) RabbitMQConsumerConfigurationBuilder
	WithQueueName(queueName string) RabbitMQConsumerConfigurationBuilder
	WithDurable(durable bool) RabbitMQConsumerConfigurationBuilder
	WithAutoDeleteQueue(autoDelete bool) RabbitMQConsumerConfigurationBuilder
	WithExclusiveQueue(exclusive bool) RabbitMQConsumerConfigurationBuilder
	WithQueueArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithExchangeName(exchangeName string) RabbitMQConsumerConfigurationBuilder
	WithAutoDeleteExchange(autoDelete bool) RabbitMQConsumerConfigurationBuilder
	WithExchangeType(exchangeType types.ExchangeType) RabbitMQConsumerConfigurationBuilder
	WithExchangeArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithRoutingKey(routingKey string) RabbitMQConsumerConfigurationBuilder
	WithBindingArgs(args map[string]any) RabbitMQConsumerConfigurationBuilder
	WithName(name string) RabbitMQConsumerConfigurationBuilder
	Build() *RabbitMQConsumerConfiguration
}

// rabbitMQConsumerConfigurationBuilder is a struct that represents the rabbitmq consumer configuration builder.
type rabbitMQConsumerConfigurationBuilder struct {
	rabbitmqConsumerConfigurations *RabbitMQConsumerConfiguration
	pipelinesBuilder               pipeline.ConsumerPipelineConfigurationBuilder
	handlersBuilder                messageConsumer.ConsumerHandlerConfigurationBuilder
}

// NewRabbitMQConsumerConfigurationBuilder creates a new rabbitmq consumer configuration builder.
func NewRabbitMQConsumerConfigurationBuilder(
	messageType types2.IMessage,
) RabbitMQConsumerConfigurationBuilder {
	return &rabbitMQConsumerConfigurationBuilder{
		rabbitmqConsumerConfigurations: NewDefaultRabbitMQConsumerConfiguration(messageType),
	}
}

// WithPipelines adds a pipeline to the rabbitmq consumer configuration.
func (b *rabbitMQConsumerConfigurationBuilder) WIthPipelines(
	pipelineBuilderFunc pipeline.ConsumerPipelineConfigurationBuilderFunc,
) RabbitMQConsumerConfigurationBuilder {
	builder := pipeline.NewConsumerPipelineConfigurationBuilder()
	if pipelineBuilderFunc != nil {
		pipelineBuilderFunc(builder)
	}
	b.pipelinesBuilder = builder

	return b
}

// WithHandlers adds a handler to the rabbitmq consumer configuration.
func (b *rabbitMQConsumerConfigurationBuilder) WithHandlers(
	consumerBuilderFunc messageConsumer.ConsumerHandlerConfigurationBuilderFunc,
) RabbitMQConsumerConfigurationBuilder {
	builder := messageConsumer.NewConsumerHandlersConfigurationBuilder()
	if consumerBuilderFunc != nil {
		consumerBuilderFunc(builder)
	}
	b.handlersBuilder = builder

	return b
}

// WithExitOnError sets the exit on error flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithExitOnError(
	exitOnError bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExitOnError = exitOnError

	return b
}

// WithName sets the name of the rabbitmq consumer configuration.
func (b *rabbitMQConsumerConfigurationBuilder) WithName(
	name string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.Name = name

	return b
}

// WithAutoAck sets the auto ack flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithAutoAck(
	ack bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.AutoAck = ack

	return b
}

// WithNoLocal sets the no local flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithNoLocal(
	noLocal bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.NoLocal = noLocal

	return b
}

// WithNoWait sets the no wait flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithNoWait(
	noWait bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.NoWait = noWait

	return b
}

// WithConcurrencyLimit sets the concurrency limit.
func (b *rabbitMQConsumerConfigurationBuilder) WithConcurrencyLimit(
	limit int,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ConcurrencyLimit = limit

	return b
}

// WithPrefetchCount sets the prefetch count.
func (b *rabbitMQConsumerConfigurationBuilder) WithPrefetchCount(
	count int,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.PrefetchCount = count

	return b
}

// WithConsumerID sets the consumer id.
func (b *rabbitMQConsumerConfigurationBuilder) WithConsumerID(
	consumerId string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ConsumerId = consumerId

	return b
}

// WithQueueName sets the queue name.
func (b *rabbitMQConsumerConfigurationBuilder) WithQueueName(
	queueName string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Name = queueName

	return b
}

func (b *rabbitMQConsumerConfigurationBuilder) WithDurable(
	durable bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Durable = durable
	b.rabbitmqConsumerConfigurations.QueueOptions.Durable = durable

	return b
}

// WithAutoDeleteQueue sets the auto delete queue flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithAutoDeleteQueue(
	autoDelete bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.AutoDelete = autoDelete

	return b
}

// WithExclusiveQueue sets the exclusive queue flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithExclusiveQueue(
	exclusive bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Exclusive = exclusive

	return b
}

// WithQueueArgs sets the queue args.
func (b *rabbitMQConsumerConfigurationBuilder) WithQueueArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.QueueOptions.Args = args

	return b
}

// WithExchangeName sets the exchange name.
func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeName(
	exchangeName string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Name = exchangeName

	return b
}

// WithAutoDeleteExchange sets the auto delete exchange flag.
func (b *rabbitMQConsumerConfigurationBuilder) WithAutoDeleteExchange(
	autoDelete bool,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.AutoDelete = autoDelete

	return b
}

// WithExchangeType sets the exchange type.
func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeType(
	exchangeType types.ExchangeType,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Type = exchangeType

	return b
}

// WithExchangeArgs sets the exchange args.
func (b *rabbitMQConsumerConfigurationBuilder) WithExchangeArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.ExchangeOptions.Args = args

	return b
}

// WithRoutingKey sets the routing key.
func (b *rabbitMQConsumerConfigurationBuilder) WithRoutingKey(
	routingKey string,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.BindingOptions.RoutingKey = routingKey

	return b
}

// WithBindingArgs sets the binding args.
func (b *rabbitMQConsumerConfigurationBuilder) WithBindingArgs(
	args map[string]any,
) RabbitMQConsumerConfigurationBuilder {
	b.rabbitmqConsumerConfigurations.BindingOptions.Args = args

	return b
}

// Build builds the rabbitmq consumer configuration.
func (b *rabbitMQConsumerConfigurationBuilder) Build() *RabbitMQConsumerConfiguration {
	if b.pipelinesBuilder != nil {
		b.rabbitmqConsumerConfigurations.Pipelines = b.pipelinesBuilder.Build().Pipelines
	}
	if b.handlersBuilder != nil {
		b.rabbitmqConsumerConfigurations.Handlers = b.handlersBuilder.Build().Handlers
	}

	return b.rabbitmqConsumerConfigurations
}
