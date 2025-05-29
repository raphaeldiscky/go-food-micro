package tracing

import (
	"go.opentelemetry.io/otel/trace"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
)

var MessagingTracer trace.Tracer

func init() {
	MessagingTracer = tracing.NewAppTracer(
		"github.com/raphaeldiscky/go-food-micro/internal/pkg/messaging",
	) // instrumentation name
}
