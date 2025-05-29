// Package producer provides a producer builder.
package producer

// ProducerBuilderFuc is a function that builds a producer.
type ProducerBuilderFuc[T interface{}] func(builder T)
