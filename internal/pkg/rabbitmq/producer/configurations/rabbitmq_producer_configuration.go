// Package configurations provides a set of functions for the rabbitmq producer configurations.
package configurations

import (
	"reflect"

	types2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/utils"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/producer/options"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"
)

// RabbitMQProducerConfiguration is a struct that contains the rabbitmq producer configuration.
type RabbitMQProducerConfiguration struct {
	ProducerMessageType reflect.Type
	ExchangeOptions     *options.RabbitMQExchangeOptions
	RoutingKey          string
	DeliveryMode        uint8
	Priority            uint8
	AppId               string
	Expiration          string
	ReplyTo             string
	ContentEncoding     string
}

// NewDefaultRabbitMQProducerConfiguration creates a new default rabbitmq producer configuration.
func NewDefaultRabbitMQProducerConfiguration(
	messageType types2.IMessage,
) *RabbitMQProducerConfiguration {
	return &RabbitMQProducerConfiguration{
		ExchangeOptions: &options.RabbitMQExchangeOptions{
			Durable: true,
			Type:    types.ExchangeTopic,
			Name:    utils.GetTopicOrExchangeName(messageType),
		},
		DeliveryMode:        2,
		RoutingKey:          utils.GetRoutingKey(messageType),
		ProducerMessageType: utils.GetMessageBaseReflectType(messageType),
	}
}
