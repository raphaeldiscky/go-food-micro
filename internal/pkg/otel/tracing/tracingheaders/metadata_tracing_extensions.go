// Package tracingheaders provides metadata tracing extensions.
package tracingheaders

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

func GetTracingTraceId(m metadata.Metadata) string {
	return m.GetString(TraceId)
}

func GetTracingParentSpanId(m metadata.Metadata) string {
	return m.GetString(ParentSpanId)
}

func GetTracingTraceparent(m metadata.Metadata) string {
	return m.GetString(Traceparent)
}
