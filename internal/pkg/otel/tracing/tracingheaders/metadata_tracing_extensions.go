// Package tracingheaders provides metadata tracing extensions.
package tracingheaders

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

// GetTracingTraceID gets the tracing trace id.
func GetTracingTraceID(m metadata.Metadata) string {
	return m.GetString(TraceId)
}

// GetTracingParentSpanID gets the tracing parent span id.
func GetTracingParentSpanID(m metadata.Metadata) string {
	return m.GetString(ParentSpanId)
}

// GetTracingTraceParent gets the tracing trace parent.
func GetTracingTraceParent(m metadata.Metadata) string {
	return m.GetString(Traceparent)
}
