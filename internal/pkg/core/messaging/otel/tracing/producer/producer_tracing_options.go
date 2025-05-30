// Package producer provides a producer tracing options.
package producer

import "go.opentelemetry.io/otel/attribute"

// ProducerTracingOptions is a struct that represents the producer tracing options.
type ProducerTracingOptions struct {
	MessagingSystem string
	DestinationKind string
	Destination     string
	OtherAttributes []attribute.KeyValue
}
