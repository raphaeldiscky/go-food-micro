// Package oteltracing provides a otel tracing middleware.
package oteltracing

import (
	"github.com/labstack/echo/v4/middleware"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	oteltrace "go.opentelemetry.io/otel/trace"
)

// config is used to configure the mux middleware.
// Ref: https://github.com/open-telemetry/opentelemetry-go-contrib/blob/main/instrumentation/github.com/labstack/echo/otelecho/echo.go
type config struct {
	tracerProvider      oteltrace.TracerProvider
	propagators         propagation.TextMapPropagator
	skipper             middleware.Skipper
	instrumentationName string
	serviceName         string
}

// Option is the option for the otel tracing middleware.
type Option interface {
	apply(*config)
}

// optionFunc is the option function for the otel tracing middleware.
type optionFunc func(*config)

// apply applies the option to the config.
func (o optionFunc) apply(c *config) {
	o(c)
}

// defualtConfig is the default config for the otel tracing middleware.
var defualtConfig = config{
	tracerProvider:      otel.GetTracerProvider(),
	propagators:         otel.GetTextMapPropagator(),
	skipper:             middleware.DefaultSkipper,
	instrumentationName: "echo",
	serviceName:         "app",
}

// WithPropagators is the option function for the otel tracing middleware.
// specifies propagators to use for extracting
// information from the HTTP requests. If none are specified, global
// ones will be used.
func WithPropagators(propagators propagation.TextMapPropagator) Option {
	return optionFunc(func(cfg *config) {
		if propagators != nil {
			cfg.propagators = propagators
		}
	})
}

// WithTracerProvider is the option function for the otel tracing middleware.
// specifies a tracer provider to use for creating a tracer.
// If none is specified, the global provider is used.
func WithTracerProvider(provider oteltrace.TracerProvider) Option {
	return optionFunc(func(cfg *config) {
		if provider != nil {
			cfg.tracerProvider = provider
		}
	})
}

// WithSkipper is the option function for the otel tracing middleware.
// specifies a skipper for allowing requests to skip generating spans.
func WithSkipper(skipper middleware.Skipper) Option {
	return optionFunc(func(cfg *config) {
		cfg.skipper = skipper
	})
}

// WithInstrumentationName is the option function for the otel tracing middleware.
func WithInstrumentationName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.instrumentationName != "" {
			cfg.instrumentationName = v
		}
	})
}

// WithServiceName is the option function for the otel tracing middleware.
func WithServiceName(v string) Option {
	return optionFunc(func(cfg *config) {
		if cfg.serviceName != "" {
			cfg.serviceName = v
		}
	})
}
