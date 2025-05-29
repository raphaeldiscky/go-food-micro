package configurations

import (
	"github.com/samber/lo"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	consumerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/consumer/configurations"
	producerConfigurations "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/configurations"
)

type RabbitMQConfigurationBuilder interface {
	AddProducer(
		producerMessageType types.IMessage,
		producerBuilderFunc producerConfigurations.RabbitMQProducerConfigurationBuilderFuc,
	) RabbitMQConfigurationBuilder
	AddConsumer(
		consumerMessageType types.IMessage,
		consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
	) RabbitMQConfigurationBuilder
	Build() *RabbitMQConfiguration
}

type rabbitMQConfigurationBuilder struct {
	rabbitMQConfiguration *RabbitMQConfiguration
	consumerBuilders      []consumerConfigurations.RabbitMQConsumerConfigurationBuilder
	producerBuilders      []producerConfigurations.RabbitMQProducerConfigurationBuilder
}

func NewRabbitMQConfigurationBuilder() RabbitMQConfigurationBuilder {
	return &rabbitMQConfigurationBuilder{
		rabbitMQConfiguration: &RabbitMQConfiguration{},
	}
}

func (r *rabbitMQConfigurationBuilder) AddProducer(
	producerMessageType types.IMessage,
	producerBuilderFunc producerConfigurations.RabbitMQProducerConfigurationBuilderFuc,
) RabbitMQConfigurationBuilder {
	builder := producerConfigurations.NewRabbitMQProducerConfigurationBuilder(producerMessageType)
	if producerBuilderFunc != nil {
		producerBuilderFunc(builder)
	}

	r.producerBuilders = append(r.producerBuilders, builder)

	return r
}

func (r *rabbitMQConfigurationBuilder) AddConsumer(
	consumerMessageType types.IMessage,
	consumerBuilderFunc consumerConfigurations.RabbitMQConsumerConfigurationBuilderFuc,
) RabbitMQConfigurationBuilder {
	builder := consumerConfigurations.NewRabbitMQConsumerConfigurationBuilder(consumerMessageType)
	if consumerBuilderFunc != nil {
		consumerBuilderFunc(builder)
	}

	r.consumerBuilders = append(r.consumerBuilders, builder)

	return r
}

func (r *rabbitMQConfigurationBuilder) Build() *RabbitMQConfiguration {
	consumersConfigs := lo.Map(
		r.consumerBuilders,
		func(builder consumerConfigurations.RabbitMQConsumerConfigurationBuilder, index int) *consumerConfigurations.RabbitMQConsumerConfiguration {
			return builder.Build()
		},
	)

	producersConfigs := lo.Map(
		r.producerBuilders,
		func(builder producerConfigurations.RabbitMQProducerConfigurationBuilder, index int) *producerConfigurations.RabbitMQProducerConfiguration {
			return builder.Build()
		},
	)

	r.rabbitMQConfiguration.ConsumersConfigurations = consumersConfigs
	r.rabbitMQConfiguration.ProducersConfigurations = producersConfigs

	return r.rabbitMQConfiguration
}
