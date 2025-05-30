// Package consumer provides a module for the consumer.
package consumer

// ConsumerBuilderFucT is a consumer builder function.
type ConsumerBuilderFucT[T interface{}] func(builder T)
