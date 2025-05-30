// Package consumer provides the consumer options.
package consumer

// ConsumerOptions is a struct that represents the consumer options.
type ConsumerOptions struct {
	// ExitOnError is a flag that indicates if the consumer should exit on error.
	ExitOnError bool
	// ConsumerId is the id of the consumer.
	ConsumerId string
}
