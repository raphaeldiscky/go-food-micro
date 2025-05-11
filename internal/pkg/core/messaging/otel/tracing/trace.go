package tracing

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	"go.opentelemetry.io/otel/trace"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/raphaeldiscky/go-food-micro/internal/pkg/messaging",
	) // instrumentation name
}
