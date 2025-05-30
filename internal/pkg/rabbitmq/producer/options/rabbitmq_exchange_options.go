// Package options provides a set of options for the rabbitmq producer.
package options

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"

// RabbitMQExchangeOptions is a struct that represents a rabbitmq exchange options.
type RabbitMQExchangeOptions struct {
	Name       string
	Type       types.ExchangeType
	AutoDelete bool
	Durable    bool
	Args       map[string]any
}
