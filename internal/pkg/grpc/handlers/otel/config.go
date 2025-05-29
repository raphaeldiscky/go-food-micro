// Package otel provides a otel config.
package otel

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// ref: https://github.com/bakins/otel-grpc-statshandler/blob/main/statshandler.go

// Option applies an option value when creating a Handler.
type Option interface {
	apply(*config)
}

// optionFunc is a function that applies an option value when creating a Handler.
type optionFunc func(*config)

func (f optionFunc) apply(c *config) {
	f(c)
}

// config is a struct that represents a config.
type config struct {
	metricsProvider     metric.MeterProvider
	tracerProvider      trace.TracerProvider
	propagator          propagation.TextMapPropagator
	Namespace           string
	serviceName         string
	instrumentationName string
}

// defualtConfig is a struct that represents a default config.
var defualtConfig = config{
	metricsProvider:     otel.GetMeterProvider(),
	tracerProvider:      otel.GetTracerProvider(),
	propagator:          otel.GetTextMapPropagator(),
	serviceName:         "application",
	instrumentationName: "grpc-otel",
}

// WithMeterProvider is a function that sets the meter provider.
func WithMeterProvider(m metric.MeterProvider) Option {
	return optionFunc(func(c *config) {
		c.metricsProvider = m
	})
}

// WithTraceProvider is a function that sets the trace provider.
func WithTraceProvider(t trace.TracerProvider) Option {
	return optionFunc(func(c *config) {
		c.tracerProvider = t
	})
}

// WithPropagators is a function that sets the propagators.
func WithPropagators(p propagation.TextMapPropagator) Option {
	return optionFunc(func(c *config) {
		c.propagator = p
	})
}

// SetNamespace is a function that sets the namespace.
func SetNamespace(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.Namespace != "" {
			cfg.Namespace = v
		}
	})
}

// SetServiceName is a function that sets the service name.
func SetServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}

// SetInstrumentationName is a function that sets the instrumentation name.
func SetInstrumentationName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.instrumentationName != "" {
			cfg.instrumentationName = v
		}
	})
}
