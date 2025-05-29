// Package bus provides a rabbitmq bus.
package bus

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"emperror.dev/errors"
	"github.com/samber/lo"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/bus"
	consumer2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/consumer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/producer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/consumercontracts"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/producercontracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/rabbitmqerrors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// RabbitmqBus is the interface for the rabbitmq bus.
type RabbitmqBus interface {
	bus.Bus
	consumerConfigurations.RabbitMQConsumerConnector
}

// rabbitmqBus is a struct that represents a rabbitmq bus.
type rabbitmqBus struct {
	messageTypeConsumers    map[reflect.Type][]consumer2.Consumer
	producer                producer.Producer
	rabbitmqConfiguration   *configurations.RabbitMQConfiguration
	rabbitmqConfigBuilder   configurations.RabbitMQConfigurationBuilder
	logger                  logger.Logger
	consumerFactory         consumercontracts.ConsumerFactory
	producerFactory         producercontracts.ProducerFactory
	isConsumedNotifications []func(message types.IMessage)
	isProducedNotifications []func(message types.IMessage)
}

// NewRabbitmqBus creates a new rabbitmq bus.
func NewRabbitmqBus(
	logger logger.Logger,
	consumerFactory consumercontracts.ConsumerFactory,
	producerFactory producercontracts.ProducerFactory,
	rabbitmqBuilderFunc configurations.RabbitMQConfigurationBuilderFuc,
) (RabbitmqBus, error) {
	builder := configurations.NewRabbitMQConfigurationBuilder()
	if rabbitmqBuilderFunc != nil {
		rabbitmqBuilderFunc(builder)
	}

	rabbitmqConfiguration := builder.Build()
	rabbitBus := &rabbitmqBus{
		logger:                logger,
		rabbitmqConfiguration: rabbitmqConfiguration,
		consumerFactory:       consumerFactory,
		producerFactory:       producerFactory,
		rabbitmqConfigBuilder: builder,
		messageTypeConsumers:  map[reflect.Type][]consumer2.Consumer{},
	}

	producersConfigurationMap := make(
		map[string]*producerConfigurations.RabbitMQProducerConfiguration,
	)
	lo.ForEach(
		rabbitBus.rabbitmqConfiguration.ProducersConfigurations,
		func(config *producerConfigurations.RabbitMQProducerConfiguration, _ int) {
			key := config.ProducerMessageType.String()
			producersConfigurationMap[key] = config
		},
	)

	consumersConfigurationMap := make(
		map[string]*consumerConfigurations.RabbitMQConsumerConfiguration,
	)
	lo.ForEach(
		rabbitBus.rabbitmqConfiguration.ConsumersConfigurations,
		func(config *consumerConfigurations.RabbitMQConsumerConfiguration, _ int) {
			key := config.ConsumerMessageType.String()
			consumersConfigurationMap[key] = config
		},
	)

	for _, consumerConfiguration := range consumersConfigurationMap {
		mqConsumer, err := consumerFactory.CreateConsumer(
			consumerConfiguration,
			// IsConsumed Notification
			func(message types.IMessage) {
				if rabbitBus.isConsumedNotifications != nil {
					for _, notification := range rabbitBus.isConsumedNotifications {
						notification(message)
					}
				}
			},
		)
		if err != nil {
			return nil, err
		}
		rabbitBus.messageTypeConsumers[consumerConfiguration.ConsumerMessageType] = append(
			rabbitBus.messageTypeConsumers[consumerConfiguration.ConsumerMessageType],
			mqConsumer,
		)
	}

	mqProducer, err := producerFactory.CreateProducer(
		producersConfigurationMap,
		// IsProduced Notification
		func(message types.IMessage) {
			if rabbitBus.isProducedNotifications != nil {
				for _, notification := range rabbitBus.isProducedNotifications {
					notification(message)
				}
			}
		},
	)
	if err != nil {
		return nil, err
	}
	rabbitBus.producer = mqProducer

	return rabbitBus, nil
}

// IsConsumed adds a notification to the rabbitmq bus.
func (r *rabbitmqBus) IsConsumed(h func(message types.IMessage)) {
	r.isConsumedNotifications = append(r.isConsumedNotifications, h)
}

// IsProduced adds a notification to the rabbitmq bus.
func (r *rabbitmqBus) IsProduced(h func(message types.IMessage)) {
	r.isProducedNotifications = append(r.isProducedNotifications, h)
}

// ConnectConsumer adds a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type.
func (r *rabbitmqBus) ConnectConsumer(
	messageType types.IMessage,
	consumer consumer2.Consumer,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

	r.messageTypeConsumers[typeName] = append(
		r.messageTypeConsumers[typeName],
		consumer,
	)

	return nil
}

// ConnectRabbitMQConsumer adds a new consumer to existing message type consumers. if there is no consumer, will create a new consumer for the message type.
func (r *rabbitmqBus) ConnectRabbitMQConsumer(
	messageType types.IMessage,
	consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

	builder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(
		messageType,
	)
	if consumerBuilderFunc != nil {
		consumerBuilderFunc(builder)
	}
	consumerConfig := builder.Build()
	mqConsumer, err := r.consumerFactory.CreateConsumer(
		consumerConfig,
		// IsConsumed Notification
		func(message types.IMessage) {
			if len(r.isConsumedNotifications) > 0 {
				for _, notification := range r.isConsumedNotifications {
					if notification != nil {
						notification(message)
					}
				}
			}
		},
	)
	if err != nil {
		return err
	}

	r.messageTypeConsumers[typeName] = append(
		r.messageTypeConsumers[typeName],
		mqConsumer,
	)

	return nil
}

// createNewConsumer creates a new consumer for the given message type and handler.
func (r *rabbitmqBus) createNewConsumer(
	messageType types.IMessage,
	consumerHandler consumer2.ConsumerHandler,
) (consumer2.Consumer, error) {
	consumerBuilder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(messageType)
	consumerBuilder.WithHandlers(func(builder consumer2.ConsumerHandlerConfigurationBuilder) {
		builder.AddHandler(consumerHandler)
	})
	consumerConfig := consumerBuilder.Build()

	return r.consumerFactory.CreateConsumer(
		consumerConfig,
		func(message types.IMessage) {
			if len(r.isConsumedNotifications) > 0 {
				for _, notification := range r.isConsumedNotifications {
					if notification != nil {
						notification(message)
					}
				}
			}
		},
	)
}

// ConnectConsumerHandler adds a handler to existing consumer. creates new consumer if not exist.
func (r *rabbitmqBus) ConnectConsumerHandler(
	messageType types.IMessage,
	consumerHandler consumer2.ConsumerHandler,
) error {
	typeName := utils.GetMessageBaseReflectType(messageType)

	consumersForType := r.messageTypeConsumers[typeName]
	if consumersForType != nil {
		for _, c := range consumersForType {
			c.ConnectHandler(consumerHandler)
		}

		return nil
	}

	mqConsumer, err := r.createNewConsumer(messageType, consumerHandler)
	if err != nil {
		return err
	}

	r.messageTypeConsumers[typeName] = append(r.messageTypeConsumers[typeName], mqConsumer)

	return nil
}

// Start starts the rabbitmq bus.
func (r *rabbitmqBus) Start(ctx context.Context) error {
	r.logger.Infof(
		"rabbitmq is running on host: %s",
		r.consumerFactory.Connection().Raw().LocalAddr().String(),
	)

	for messageType, consumers := range r.messageTypeConsumers {
		name := typeMapper.GetTypeNameByType(messageType)
		r.logger.Info(fmt.Sprintf("consuming message type %s", name))
		for _, rabbitConsumer := range consumers {
			err := rabbitConsumer.Start(ctx)
			r.logger.Info(
				fmt.Sprintf("consumer %s, started", rabbitConsumer.GetName()),
			)
			if errors.Is(err, rabbitmqerrors.ErrDisconnected) {
				r.logger.Info(
					fmt.Sprintf(
						"consumer %s, disconnected with err: %v",
						rabbitConsumer.GetName(),
						err,
					),
				)
				// will process again with reConsume functionality
				continue
			} else if err != nil {
				r.logger.Error(
					fmt.Sprintf(
						"error in consumer %s, with err: %v",
						rabbitConsumer.GetName(),
						err,
					),
				)
				err2 := r.Stop()
				if err2 != nil {
					return errors.WrapIf(err, err2.Error())
				}

				return err
			}
		}
	}

	return nil
}

// Stop stops the rabbitmq bus.
func (r *rabbitmqBus) Stop() error {
	waitGroup := sync.WaitGroup{}

	for _, consumers := range r.messageTypeConsumers {
		for _, c := range consumers {
			waitGroup.Add(1)

			go func(c consumer2.Consumer) {
				defer waitGroup.Done()

				err := c.Stop()
				if err != nil {
					r.logger.Error("error in the unconsuming")
				}
			}(c)
		}
	}
	waitGroup.Wait()

	// err := r.rabbitmqConnection.Close()
	// if err == amqp091.ErrClosed {
	//	return nil
	//}

	return nil
}

// PublishMessage publishes a message to the rabbitmq bus.
func (r *rabbitmqBus) PublishMessage(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
) error {
	if r.producer == nil {
		r.logger.Fatal("can't find a producer for publishing messages")
	}

	return r.producer.PublishMessage(ctx, message, meta)
}

// PublishMessageWithTopicName publishes a message to the rabbitmq bus with a topic name.
func (r *rabbitmqBus) PublishMessageWithTopicName(
	ctx context.Context,
	message types.IMessage,
	meta metadata.Metadata,
	topicOrExchangeName string,
) error {
	return r.producer.PublishMessageWithTopicName(
		ctx,
		message,
		meta,
		topicOrExchangeName,
	)
}
