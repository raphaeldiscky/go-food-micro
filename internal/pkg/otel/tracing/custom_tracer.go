// Package tracing provides a custom tracer.
package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
)

// AppTracer is a custom tracer.
type AppTracer interface {
	trace.Tracer
}

// appTracer is a custom tracer.
type appTracer struct {
	trace.Tracer
}

// Start starts a new span.
func (c *appTracer) Start(
	ctx context.Context,
	spanName string,
	opts ...trace.SpanStartOption,
) (context.Context, trace.Span) {
	parentSpan := trace.SpanFromContext(ctx)
	if parentSpan != nil {
		utils.ContextWithParentSpan(ctx, parentSpan)
	}

	return c.Tracer.Start(ctx, spanName, opts...)
}

// NewAppTracer creates a new custom tracer.
func NewAppTracer(name string, options ...trace.TracerOption) AppTracer {
	// without registering `NewOtelTracing` it uses global empty (NoopTracer) TraceProvider but after using `NewOtelTracing`, global TraceProvider will be replaced
	tracer := otel.Tracer(name, options...)

	return &appTracer{Tracer: tracer}
}
