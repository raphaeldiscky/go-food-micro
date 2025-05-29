// Package configurations provides a set of functions for the rabbitmq producer configurations.
package configurations

import (
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// RabbitMQProducerConfigurationBuilder is a interface that contains the rabbitmq producer configuration builder.
type RabbitMQProducerConfigurationBuilder interface {
	WithDurable(durable bool) RabbitMQProducerConfigurationBuilder
	WithAutoDeleteExchange(autoDelete bool) RabbitMQProducerConfigurationBuilder
	WithExchangeType(exchangeType types.ExchangeType) RabbitMQProducerConfigurationBuilder
	WithExchangeName(exchangeName string) RabbitMQProducerConfigurationBuilder
	WithRoutingKey(routingKey string) RabbitMQProducerConfigurationBuilder
	WithExchangeArgs(args map[string]any) RabbitMQProducerConfigurationBuilder
	WithDeliveryMode(deliveryMode uint8) RabbitMQProducerConfigurationBuilder
	WithPriority(priority uint8) RabbitMQProducerConfigurationBuilder
	WithAppId(appId string) RabbitMQProducerConfigurationBuilder
	WithExpiration(expiration string) RabbitMQProducerConfigurationBuilder
	WithReplyTo(replyTo string) RabbitMQProducerConfigurationBuilder
	WithContentEncoding(contentEncoding string) RabbitMQProducerConfigurationBuilder
	Build() *RabbitMQProducerConfiguration
}

// rabbitMQProducerConfigurationBuilder is a struct that contains the rabbitmq producer configuration builder.
type rabbitMQProducerConfigurationBuilder struct {
	rabbitmqProducerOptions *RabbitMQProducerConfiguration
}

// NewRabbitMQProducerConfigurationBuilder creates a new rabbitmq producer configuration builder.
func NewRabbitMQProducerConfigurationBuilder(
	messageType types2.IMessage,
) RabbitMQProducerConfigurationBuilder {
	return &rabbitMQProducerConfigurationBuilder{
		rabbitmqProducerOptions: NewDefaultRabbitMQProducerConfiguration(messageType),
	}
}

// WithDurable sets the durable option.
func (b *rabbitMQProducerConfigurationBuilder) WithDurable(
	durable bool,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Durable = durable

	return b
}

// WithAutoDeleteExchange sets the auto delete exchange option.
func (b *rabbitMQProducerConfigurationBuilder) WithAutoDeleteExchange(
	autoDelete bool,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.AutoDelete = autoDelete

	return b
}

// WithExchangeType sets the exchange type option.
func (b *rabbitMQProducerConfigurationBuilder) WithExchangeType(
	exchangeType types.ExchangeType,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Type = exchangeType

	return b
}

// WithRoutingKey sets the routing key option.
func (b *rabbitMQProducerConfigurationBuilder) WithRoutingKey(
	routingKey string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.RoutingKey = routingKey

	return b
}

// WithExchangeName sets the exchange name option.
func (b *rabbitMQProducerConfigurationBuilder) WithExchangeName(
	exchangeName string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Name = exchangeName

	return b
}

// WithExchangeArgs sets the exchange args option.
func (b *rabbitMQProducerConfigurationBuilder) WithExchangeArgs(
	args map[string]any,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ExchangeOptions.Args = args

	return b
}

// WithDeliveryMode sets the delivery mode option.
func (b *rabbitMQProducerConfigurationBuilder) WithDeliveryMode(
	deliveryMode uint8,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.DeliveryMode = deliveryMode

	return b
}

// WithPriority sets the priority option.
func (b *rabbitMQProducerConfigurationBuilder) WithPriority(
	priority uint8,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Priority = priority

	return b
}

// WithAppId sets the app id option.
func (b *rabbitMQProducerConfigurationBuilder) WithAppId(
	appId string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.AppId = appId

	return b
}

// WithExpiration sets the expiration option.
func (b *rabbitMQProducerConfigurationBuilder) WithExpiration(
	expiration string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.Expiration = expiration

	return b
}

func (b *rabbitMQProducerConfigurationBuilder) WithReplyTo(
	replyTo string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ReplyTo = replyTo

	return b
}

// WithContentEncoding sets the content encoding option.
func (b *rabbitMQProducerConfigurationBuilder) WithContentEncoding(
	contentEncoding string,
) RabbitMQProducerConfigurationBuilder {
	b.rabbitmqProducerOptions.ContentEncoding = contentEncoding

	return b
}

// Build builds the rabbitmq producer configuration.
func (b *rabbitMQProducerConfigurationBuilder) Build() *RabbitMQProducerConfiguration {
	return b.rabbitmqProducerOptions
}
