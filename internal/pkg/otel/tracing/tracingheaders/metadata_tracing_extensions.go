// Package tracingheaders provides metadata tracing extensions.
package tracingheaders

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

// GetTracingTraceId gets the tracing trace id.
func GetTracingTraceId(m metadata.Metadata) string {
	return m.GetString(TraceId)
}

// GetTracingParentSpanId gets the tracing parent span id.
func GetTracingParentSpanId(m metadata.Metadata) string {
	return m.GetString(ParentSpanId)
}

// GetTracingTraceparent gets the tracing traceparent.
func GetTracingTraceparent(m metadata.Metadata) string {
	return m.GetString(Traceparent)
}
