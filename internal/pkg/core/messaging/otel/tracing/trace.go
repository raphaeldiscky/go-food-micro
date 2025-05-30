// Package tracing provides a tracing.
package tracing

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
)

// MessagingTracer is a tracer for the messaging system.
var MessagingTracer trace.Tracer

// init is a function that initializes the tracing.
//
//nolint:gochecknoinits // This is a standard pattern for initializing the tracing
func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/raphaeldiscky/go-food-micro/internal/pkg/messaging",
	) // instrumentation name
}
