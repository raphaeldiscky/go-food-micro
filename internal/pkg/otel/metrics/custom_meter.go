// Package metrics provides a custom meter.
package metrics

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

// AppMetrics is a custom meter.
type AppMetrics interface {
	metric.Meter
}

// appMetrics is a custom meter.
type appMetrics struct {
	metric.Meter
}

// NewAppMeter creates a new custom meter.
func NewAppMeter(name string, opts ...metric.MeterOption) AppMetrics {
	// Meter can be a global/package variable.
	// https://github.com/open-telemetry/opentelemetry-go/blob/46f2ce5ca6adaa264c37cdbba251c9184a06ed7f/metric.go#LL35C6-L35C11
	meter := otel.Meter(name, opts...)

	return &appMetrics{Meter: meter}
}
