package options

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"

type RabbitMQExchangeOptions struct {
	Name       string
	Type       types.ExchangeType
	AutoDelete bool
	Durable    bool
	Args       map[string]any
}
