// Package options provides a set of functions for the rabbitmq consumer options.
package options

// RabbitMQQueueOptions is a struct that contains the rabbitmq queue options.
type RabbitMQQueueOptions struct {
	Name       string
	Durable    bool
	Exclusive  bool
	AutoDelete bool
	Args       map[string]any
}
