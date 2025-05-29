// Package options provides a set of functions for the rabbitmq consumer options.
package options

// RabbitMQBindingOptions is a struct that contains the rabbitmq binding options.
type RabbitMQBindingOptions struct {
	RoutingKey string
	Args       map[string]any
}
