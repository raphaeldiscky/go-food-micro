// Package consumer provides a set of functions for the rabbitmq consumer.
package consumer

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	linq "github.com/ahmetb/go-linq/v3"
	retry "github.com/avast/retry-go"
	amqp091 "github.com/rabbitmq/amqp091-go"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	consumertracing "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/otel/tracing/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/pipeline"
	messagingTypes "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/rabbitmqerrors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
	errorutils "github.com/raphaeldiscky/go-food-micro/internal/pkg/utils/errorutils"
)

const (
	retryAttempts = 3
	retryDelay    = 300 * time.Millisecond
)

var retryOptions = []retry.Option{
	retry.Attempts(retryAttempts),
	retry.Delay(retryDelay),
	retry.DelayType(retry.BackOffDelay),
}

// rabbitMQConsumer is a struct that contains the rabbitmq consumer.
type rabbitMQConsumer struct {
	rabbitmqConsumerOptions *configurations.RabbitMQConsumerConfiguration
	connection              types.IConnection
	channel                 *amqp091.Channel
	deliveryRoutines        chan struct{} // chan should init before using channel
	messageSerializer       serializer.MessageSerializer
	logger                  logger.Logger
	rabbitmqOptions         *config.RabbitmqOptions
	ErrChan                 chan error
	handlers                []consumer.ConsumerHandler
	pipelines               []pipeline.ConsumerPipeline
	isConsumedNotifications []func(message messagingTypes.IMessage)
}

// NewRabbitMQConsumer creates a new generic RabbitMQ consumer.
func NewRabbitMQConsumer(
	rabbitmqOptions *config.RabbitmqOptions,
	connection types.IConnection,
	consumerConfiguration *configurations.RabbitMQConsumerConfiguration,
	messageSerializer serializer.MessageSerializer,
	logger logger.Logger,
	isConsumedNotifications ...func(message messagingTypes.IMessage),
) (consumer.Consumer, error) {
	if consumerConfiguration == nil {
		return nil, errors.New("consumer configuration is required")
	}

	if consumerConfiguration.ConsumerMessageType == nil {
		return nil, errors.New(
			"consumer ConsumerMessageType property is required",
		)
	}

	deliveryRoutines := make(
		chan struct{},
		consumerConfiguration.ConcurrencyLimit,
	)
	cons := &rabbitMQConsumer{
		messageSerializer:       messageSerializer,
		rabbitmqOptions:         rabbitmqOptions,
		logger:                  logger,
		rabbitmqConsumerOptions: consumerConfiguration,
		deliveryRoutines:        deliveryRoutines,
		ErrChan:                 make(chan error),
		connection:              connection,
		handlers:                consumerConfiguration.Handlers,
		pipelines:               consumerConfiguration.Pipelines,
	}

	cons.isConsumedNotifications = isConsumedNotifications

	return cons, nil
}

// IsConsumed adds a new consumed notification.
func (r *rabbitMQConsumer) IsConsumed(h func(message messagingTypes.IMessage)) {
	r.isConsumedNotifications = append(r.isConsumedNotifications, h)
}

// getExchangeName returns the exchange name for the consumer.
func (r *rabbitMQConsumer) getExchangeName() string {
	if r.rabbitmqConsumerOptions.ExchangeOptions.Name != "" {
		return r.rabbitmqConsumerOptions.ExchangeOptions.Name
	}

	return utils.GetTopicOrExchangeNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
}

// getRoutingKey returns the routing key for the consumer.
func (r *rabbitMQConsumer) getRoutingKey() string {
	if r.rabbitmqConsumerOptions.BindingOptions.RoutingKey != "" {
		return r.rabbitmqConsumerOptions.BindingOptions.RoutingKey
	}

	return utils.GetRoutingKeyFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
}

// getQueueName returns the queue name for the consumer.
func (r *rabbitMQConsumer) getQueueName() string {
	if r.rabbitmqConsumerOptions.QueueOptions.Name != "" {
		return r.rabbitmqConsumerOptions.QueueOptions.Name
	}

	return utils.GetQueueNameFromType(r.rabbitmqConsumerOptions.ConsumerMessageType)
}

// setupChannel sets up the channel for the consumer.
func (r *rabbitMQConsumer) setupChannel() error {
	var err error
	r.channel, err = r.connection.Channel()
	if err != nil {
		return rabbitmqerrors.ErrDisconnected
	}

	prefetchCount := r.rabbitmqConsumerOptions.ConcurrencyLimit * r.rabbitmqConsumerOptions.PrefetchCount

	return r.channel.Qos(prefetchCount, 0, false)
}

// setupExchange sets up the exchange for the consumer.
func (r *rabbitMQConsumer) setupExchange(exchange string) error {
	return r.channel.ExchangeDeclare(
		exchange,
		string(r.rabbitmqConsumerOptions.ExchangeOptions.Type),
		r.rabbitmqConsumerOptions.ExchangeOptions.Durable,
		r.rabbitmqConsumerOptions.ExchangeOptions.AutoDelete,
		false,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.ExchangeOptions.Args,
	)
}

// handleMessages handles messages from the channel.
func (r *rabbitMQConsumer) handleMessages(
	ctx context.Context,
	msgs <-chan amqp091.Delivery,
	chClosedCh chan *amqp091.Error,
) {
	for {
		select {
		case <-ctx.Done():
			r.logger.Info("shutting down consumer")

			return
		case amqErr := <-chClosedCh:
			// This case handles the event of closed channel e.g. abnormal shutdown
			r.logger.Errorf("AMQP Channel closed due to: %s", amqErr)

			// Re-set channel to receive notifications
			chClosedCh = make(chan *amqp091.Error, 1)
			r.channel.NotifyClose(chClosedCh)
		case msg, ok := <-msgs:
			if !ok {
				r.logger.Info("consumer connection dropped")

				return
			}

			// handle received message and remove message form queue with a manual ack
			r.handleReceived(ctx, msg)
		}
	}
}

// Start starts the rabbitmq consumer.
func (r *rabbitMQConsumer) Start(ctx context.Context) error {
	if r.connection == nil {
		return errors.New("connection is nil")
	}

	exchange := r.getExchangeName()
	routingKey := r.getRoutingKey()
	queue := r.getQueueName()

	r.reConsumeOnDropConnection(ctx)

	if err := r.setupChannel(); err != nil {
		return err
	}

	if err := r.setupExchange(exchange); err != nil {
		return err
	}

	_, err := r.channel.QueueDeclare(
		queue,
		r.rabbitmqConsumerOptions.QueueOptions.Durable,
		r.rabbitmqConsumerOptions.QueueOptions.AutoDelete,
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.QueueOptions.Args)
	if err != nil {
		return err
	}

	err = r.channel.QueueBind(
		queue,
		routingKey,
		exchange,
		r.rabbitmqConsumerOptions.NoWait,
		r.rabbitmqConsumerOptions.BindingOptions.Args)
	if err != nil {
		return err
	}

	msgs, err := r.channel.Consume(
		queue,
		r.rabbitmqConsumerOptions.ConsumerId,
		r.rabbitmqConsumerOptions.AutoAck,
		r.rabbitmqConsumerOptions.QueueOptions.Exclusive,
		r.rabbitmqConsumerOptions.NoLocal,
		r.rabbitmqConsumerOptions.NoWait,
		nil,
	)
	if err != nil {
		return err
	}

	// This channel will receive a notification when a channel closed event happens.
	chClosedCh := make(chan *amqp091.Error, 1)
	r.channel.NotifyClose(chClosedCh)

	for i := 0; i < r.rabbitmqConsumerOptions.ConcurrencyLimit; i++ {
		r.logger.Infof("Processing messages on thread %d", i)
		go r.handleMessages(ctx, msgs, chClosedCh)
	}

	return nil
}

// Stop stops the rabbitmq consumer.
func (r *rabbitMQConsumer) Stop() error {
	defer func() {
		if r.channel != nil && !r.channel.IsClosed() {
			if err := r.channel.Cancel(r.rabbitmqConsumerOptions.ConsumerId, false); err != nil {
				r.logger.Error(
					"error in canceling consumer: %v",
					err,
				)
			}
			if err := r.channel.Close(); err != nil {
				r.logger.Error(
					"error in closing channel: %v",
					err,
				)
			}
		}
	}()

	done := make(chan struct{}, 1)

	go func() {
		for {
			if len(r.deliveryRoutines) == 0 {
				done <- struct{}{}

				return
			}
		}
	}()

	<-done

	return nil
}

// ConnectHandler adds a new consumer handler.
func (r *rabbitMQConsumer) ConnectHandler(handler consumer.ConsumerHandler) {
	r.handlers = append(r.handlers, handler)
}

// GetName returns the name of the rabbitmq consumer.
func (r *rabbitMQConsumer) GetName() string {
	return r.rabbitmqConsumerOptions.Name
}

// reConsumeOnDropConnection reconnects the rabbitmq consumer on drop connection.
func (r *rabbitMQConsumer) reConsumeOnDropConnection(ctx context.Context) {
	go func() {
		defer errorutils.HandlePanic()
		for reconnect := range r.connection.ReconnectedChannel() {
			if !reflect.ValueOf(reconnect).IsValid() {
				continue
			}

			r.logger.Info("reconsume_on_drop_connection started")
			err := r.Start(ctx)
			if err != nil {
				r.logger.Error(
					"reconsume_on_drop_connection finished with error: %v",
					err,
				)

				continue
			}
			r.logger.Info(
				"reconsume_on_drop_connection finished successfully",
			)

			return
		}
	}()
}

func (r *rabbitMQConsumer) createAckFunc(
	delivery amqp091.Delivery,
	beforeConsumeSpan trace.Span,
) func() {
	return func() {
		if err := delivery.Ack(false); err != nil {
			r.logger.Error(
				"error sending ACK to RabbitMQ consumer: %v",
				consumertracing.FinishConsumerSpan(beforeConsumeSpan, err),
			)

			return
		}
		r.finishSpanAndNotify(beforeConsumeSpan, nil, delivery)
	}
}

func (r *rabbitMQConsumer) createNackFunc(
	delivery amqp091.Delivery,
	beforeConsumeSpan trace.Span,
) func() {
	return func() {
		if err := delivery.Nack(false, true); err != nil {
			r.logger.Error(
				"error in sending Nack to RabbitMQ consumer: %v",
				consumertracing.FinishConsumerSpan(beforeConsumeSpan, err),
			)

			return
		}
		r.finishSpanAndNotify(beforeConsumeSpan, nil, delivery)
	}
}

func (r *rabbitMQConsumer) finishSpanAndNotify(
	span trace.Span,
	err error,
	delivery amqp091.Delivery,
) {
	if err := consumertracing.FinishConsumerSpan(span, err); err != nil {
		r.logger.Error("error in finishing consumer span: %v", err)
	}
	if len(r.isConsumedNotifications) > 0 {
		for _, notification := range r.isConsumedNotifications {
			if notification != nil {
				notification(r.createConsumeContext(delivery).Message())
			}
		}
	}
}

// handleReceived handles the received message.
func (r *rabbitMQConsumer) handleReceived(ctx context.Context, delivery amqp091.Delivery) {
	r.deliveryRoutines <- struct{}{}
	defer func() { <-r.deliveryRoutines }()

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumerTraceOption := &consumertracing.ConsumerTracingOptions{
		MessagingSystem: "rabbitmq",
		DestinationKind: "queue",
		Destination:     r.rabbitmqConsumerOptions.QueueOptions.Name,
		OtherAttributes: []attribute.KeyValue{
			semconv.MessagingRabbitmqDestinationRoutingKey(delivery.RoutingKey),
		},
	}
	ctx, beforeConsumeSpan := consumertracing.StartConsumerSpan(
		ctx,
		&meta,
		string(delivery.Body),
		consumerTraceOption,
	)
	consumeContext := r.createConsumeContext(delivery)

	var ack, nack func()
	if !r.rabbitmqConsumerOptions.AutoAck {
		ack = r.createAckFunc(delivery, beforeConsumeSpan)
		nack = r.createNackFunc(delivery, beforeConsumeSpan)
	}

	r.handle(ctx, ack, nack, consumeContext)
}

// handle handles the message.
func (r *rabbitMQConsumer) handle(
	ctx context.Context,
	ack func(),
	nack func(),
	messageConsumeContext messagingTypes.MessageConsumeContext,
) {
	var err error
	for _, handler := range r.handlers {
		err = r.runHandlersWithRetry(ctx, handler, messageConsumeContext)
		if err != nil {
			break
		}
	}

	if err != nil {
		r.logger.Error(
			"[rabbitMQConsumer.Handle] error in handling consume message of RabbitmqMQ, prepare for nacking message",
		)
		if nack != nil && !r.rabbitmqConsumerOptions.AutoAck {
			nack()
		}
	} else if err == nil && ack != nil && !r.rabbitmqConsumerOptions.AutoAck {
		ack()
	}
}

func (r *rabbitMQConsumer) runHandlersWithRetry(
	ctx context.Context,
	handler consumer.ConsumerHandler,
	messageConsumeContext messagingTypes.MessageConsumeContext,
) error {
	err := retry.Do(func() error {
		var lastHandler pipeline.ConsumerHandlerFunc

		if len(r.pipelines) > 0 {
			reversPipes := r.reversOrder(r.pipelines)
			lastHandler = func(ctx context.Context) error {
				return handler.Handle(ctx, messageConsumeContext)
			}

			aggregateResult := linq.From(reversPipes).
				AggregateWithSeedT(lastHandler, func(next pipeline.ConsumerHandlerFunc, pipe pipeline.ConsumerPipeline) pipeline.ConsumerHandlerFunc {
					pipeValue := pipe
					nexValue := next

					return func(ctx context.Context) error {
						return pipeValue.Handle(
							ctx,
							messageConsumeContext,
							nexValue,
						)
					}
				})

			v, ok := aggregateResult.(pipeline.ConsumerHandlerFunc)
			if !ok {
				return errors.New(
					"failed to convert aggregateResult to pipeline.ConsumerHandlerFunc",
				)
			}
			err := v(ctx)
			if err != nil {
				return errors.Wrap(
					err,
					"error handling consumer handlers pipeline",
				)
			}

			return nil
		}
		err := handler.Handle(ctx, messageConsumeContext)
		if err != nil {
			return err
		}

		return nil
	}, append(retryOptions, retry.Context(ctx))...)

	return err
}

func (r *rabbitMQConsumer) createConsumeContext(
	delivery amqp091.Delivery,
) messagingTypes.MessageConsumeContext {
	message := r.deserializeData(
		delivery.ContentType,
		delivery.Type,
		delivery.Body,
	)

	var meta metadata.Metadata
	if delivery.Headers != nil {
		meta = metadata.MapToMetadata(delivery.Headers)
	}

	consumeContext := messagingTypes.NewMessageConsumeContext(
		message,
		meta,
		delivery.ContentType,
		delivery.Type,
		delivery.Timestamp,
		delivery.DeliveryTag,
		delivery.MessageId,
		delivery.CorrelationId,
	)

	return consumeContext
}

func (r *rabbitMQConsumer) deserializeData(
	contentType string,
	eventType string,
	body []byte,
) messagingTypes.IMessage {
	if contentType == "" {
		contentType = "application/json"
	}

	if len(body) == 0 {
		r.logger.Error("message body is nil or empty in the consumer")

		return nil
	}

	if contentType == "application/json" {
		// r.rabbitmqConsumerOptions.ConsumerMessageType --> actual type
		// deserialize, err := r.messageSerializer.DeserializeType(body, r.rabbitmqConsumerOptions.ConsumerMessageType, contentType)
		deserialize, err := r.messageSerializer.Deserialize(
			body,
			eventType,
			contentType,
		) // or this to explicit type deserialization
		if err != nil {
			r.logger.Errorf(
				fmt.Sprintf(
					"error in deserilizng of type '%s' in the consumer",
					eventType,
				),
			)

			return nil
		}

		return deserialize
	}

	return nil
}

func (r *rabbitMQConsumer) reversOrder(
	values []pipeline.ConsumerPipeline,
) []pipeline.ConsumerPipeline {
	var reverseValues []pipeline.ConsumerPipeline

	for i := len(values) - 1; i >= 0; i-- {
		reverseValues = append(reverseValues, values[i])
	}

	return reverseValues
}
