// Package options provides a set of functions for the rabbitmq consumer options.
package options

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/types"

// RabbitMQExchangeOptions is a struct that contains the rabbitmq exchange options.
type RabbitMQExchangeOptions struct {
	Name       string
	Type       types.ExchangeType
	AutoDelete bool
	Durable    bool
	Args       map[string]any
}
