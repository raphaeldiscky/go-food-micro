// Package consumer provides a consumer tracing options.
package consumer

import "go.opentelemetry.io/otel/attribute"

// ConsumerTracingOptions is a struct that represents the consumer tracing options.
type ConsumerTracingOptions struct {
	MessagingSystem string
	DestinationKind string
	Destination     string
	OtherAttributes []attribute.KeyValue
}
