// Package producer provides a set of functions for the rabbitmq producer.
package producer

import (
	"context"
	"time"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/attribute"

	amqp091 "github.com/rabbitmq/amqp091-go"
	uuid "github.com/satori/go.uuid"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	messageHeader "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/messageheader"
	producer3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/otel/tracing/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// rabbitMQProducer is a struct that contains the rabbitmq producer.
type rabbitMQProducer struct {
	logger                  logger.Logger
	rabbitmqOptions         *config.RabbitmqOptions
	connection              types.IConnection
	messageSerializer       serializer.MessageSerializer
	producersConfigurations map[string]*configurations.RabbitMQProducerConfiguration
	isProducedNotifications []func(message types2.IMessage)
}

// NewRabbitMQProducer creates a new rabbitmq producer.
func NewRabbitMQProducer(
	cfg *config.RabbitmqOptions,
	connection types.IConnection,
	rabbitmqProducersConfiguration map[string]*configurations.RabbitMQProducerConfiguration,
	logger logger.Logger,
	eventSerializer serializer.MessageSerializer,
	isProducedNotifications ...func(message types2.IMessage),
) (producer.Producer, error) {
	p := &rabbitMQProducer{
		logger:                  logger,
		rabbitmqOptions:         cfg,
		connection:              connection,
		messageSerializer:       eventSerializer,
		producersConfigurations: rabbitmqProducersConfiguration,
	}

	p.isProducedNotifications = isProducedNotifications

	return p, nil
}

// IsProduced adds a new produced notification.
func (r *rabbitMQProducer) IsProduced(h func(message types2.IMessage)) {
	r.isProducedNotifications = append(r.isProducedNotifications, h)
}

// PublishMessage publishes a message to the rabbitmq.
func (r *rabbitMQProducer) PublishMessage(
	ctx context.Context,
	message types2.IMessage,
	meta metadata.Metadata,
) error {
	return r.PublishMessageWithTopicName(ctx, message, meta, "")
}

// getProducerConfigurationByMessage gets the producer configuration by message.
func (r *rabbitMQProducer) getProducerConfigurationByMessage(
	message types2.IMessage,
) *configurations.RabbitMQProducerConfiguration {
	messageType := utils.GetMessageBaseReflectType(message)

	return r.producersConfigurations[messageType.String()]
}

// getExchangeAndRoutingKey determines the exchange and routing key for message publishing.
func (r *rabbitMQProducer) getExchangeAndRoutingKey(
	message types2.IMessage,
	producerConfiguration *configurations.RabbitMQProducerConfiguration,
	topicOrExchangeName string,
) (string, string) {
	var exchange string
	if topicOrExchangeName != "" {
		exchange = topicOrExchangeName
	} else if producerConfiguration != nil && producerConfiguration.ExchangeOptions.Name != "" {
		exchange = producerConfiguration.ExchangeOptions.Name
	} else {
		exchange = utils.GetTopicOrExchangeName(message)
	}

	var routingKey string
	if producerConfiguration != nil && producerConfiguration.RoutingKey != "" {
		routingKey = producerConfiguration.RoutingKey
	} else {
		routingKey = utils.GetRoutingKey(message)
	}

	return exchange, routingKey
}

// setupChannel sets up the channel for publishing with confirmation.
func (r *rabbitMQProducer) setupChannel(
	producerConfiguration *configurations.RabbitMQProducerConfiguration,
	exchange string,
) (*amqp091.Channel, chan amqp091.Confirmation, error) {
	if r.connection == nil {
		return nil, nil, errors.New("connection is nil")
	}

	if r.connection.IsClosed() {
		return nil, nil, errors.New("connection is closed, wait for connection alive")
	}

	channel, err := r.connection.Channel()
	if err != nil {
		return nil, nil, err
	}

	if err := r.ensureExchange(producerConfiguration, channel, exchange); err != nil {
		if closeErr := channel.Close(); closeErr != nil {
			r.logger.Errorf("Error closing channel after ensure exchange error: %v", closeErr)
		}

		return nil, nil, err
	}

	if err := channel.Confirm(false); err != nil {
		if closeErr := channel.Close(); closeErr != nil {
			r.logger.Errorf("Error closing channel after confirm error: %v", closeErr)
		}

		return nil, nil, err
	}

	confirms := make(chan amqp091.Confirmation)
	channel.NotifyPublish(confirms)

	return channel, confirms, nil
}

// publishMessageToChannel publishes a message to the channel and handles confirmation.
func (r *rabbitMQProducer) publishMessageToChannel(
	ctx context.Context,
	channel *amqp091.Channel,
	confirms chan amqp091.Confirmation,
	exchange string,
	routingKey string,
	props amqp091.Publishing,
) error {
	if err := channel.PublishWithContext(ctx, exchange, routingKey, true, false, props); err != nil {
		return err
	}

	if confirmed := <-confirms; !confirmed.Ack {
		return errors.New("ack not confirmed")
	}

	return nil
}

// PublishMessageWithTopicName publishes a message to the rabbitmq with topic name.
func (r *rabbitMQProducer) PublishMessageWithTopicName(
	ctx context.Context,
	message types2.IMessage,
	meta metadata.Metadata,
	topicOrExchangeName string,
) error {
	producerConfiguration := r.getProducerConfigurationByMessage(message)
	if producerConfiguration == nil {
		producerConfiguration = configurations.NewDefaultRabbitMQProducerConfiguration(message)
	}

	exchange, routingKey := r.getExchangeAndRoutingKey(
		message,
		producerConfiguration,
		topicOrExchangeName,
	)
	meta = r.getMetadata(message, meta)

	producerOptions := &producer3.ProducerTracingOptions{
		MessagingSystem: "rabbitmq",
		DestinationKind: "exchange",
		Destination:     exchange,
		OtherAttributes: []attribute.KeyValue{
			semconv.MessagingRabbitmqDestinationRoutingKey(routingKey),
		},
	}

	serializedObj, err := r.messageSerializer.Serialize(message)
	if err != nil {
		return err
	}

	ctx, beforeProduceSpan := producer3.StartProducerSpan(
		ctx,
		message,
		&meta,
		string(serializedObj.Data),
		producerOptions,
	)

	channel, confirms, err := r.setupChannel(producerConfiguration, exchange)
	if err != nil {
		return producer3.FinishProducerSpan(beforeProduceSpan, err)
	}
	defer func() {
		if err := channel.Close(); err != nil {
			r.logger.Errorf("Error closing channel: %v", err)
		}
	}()

	props := amqp091.Publishing{
		CorrelationId:   messageHeader.GetCorrelationId(meta),
		MessageId:       message.GeMessageId(),
		Timestamp:       time.Now(),
		Headers:         metadata.MetadataToMap(meta),
		Type:            message.GetMessageTypeName(),
		ContentType:     serializedObj.ContentType,
		Body:            serializedObj.Data,
		DeliveryMode:    producerConfiguration.DeliveryMode,
		Expiration:      producerConfiguration.Expiration,
		AppId:           producerConfiguration.AppId,
		Priority:        producerConfiguration.Priority,
		ReplyTo:         producerConfiguration.ReplyTo,
		ContentEncoding: producerConfiguration.ContentEncoding,
	}

	if err := r.publishMessageToChannel(ctx, channel, confirms, exchange, routingKey, props); err != nil {
		return producer3.FinishProducerSpan(beforeProduceSpan, err)
	}

	for _, notification := range r.isProducedNotifications {
		if notification != nil {
			notification(message)
		}
	}

	return producer3.FinishProducerSpan(beforeProduceSpan, nil)
}

// getMetadata gets the metadata.
func (r *rabbitMQProducer) getMetadata(
	message types2.IMessage,
	meta metadata.Metadata,
) metadata.Metadata {
	meta = metadata.FromMetadata(meta)

	// just message type name not full type name because in other side package name for type could be different
	messageHeader.SetMessageType(meta, message.GetMessageTypeName())
	messageHeader.SetMessageContentType(meta, r.messageSerializer.ContentType())

	if messageHeader.GetMessageId(meta) == "" {
		messageHeader.SetMessageId(meta, message.GeMessageId())
	}

	if messageHeader.GetMessageCreated(meta).Equal(*new(time.Time)) {
		messageHeader.SetMessageCreated(meta, message.GetCreated())
	}

	if messageHeader.GetCorrelationId(meta) == "" {
		cid := uuid.NewV4().String()
		messageHeader.SetCorrelationId(meta, cid)
	}
	messageHeader.SetMessageName(meta, utils.GetMessageName(message))

	return meta
}

// ensureExchange ensures the exchange.
func (r *rabbitMQProducer) ensureExchange(
	producersConfigurations *configurations.RabbitMQProducerConfiguration,
	channel *amqp091.Channel,
	exchangeName string,
) error {
	err := channel.ExchangeDeclare(
		exchangeName,
		string(producersConfigurations.ExchangeOptions.Type),
		producersConfigurations.ExchangeOptions.Durable,
		producersConfigurations.ExchangeOptions.AutoDelete,
		false,
		false,
		producersConfigurations.ExchangeOptions.Args,
	)
	if err != nil {
		return err
	}

	return nil
}
