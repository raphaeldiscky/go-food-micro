// Package tracing provides a module for the tracing.
package tracing

import (
	"context"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Module is a module for the tracing.
var (
	// Module provided to fxlog
	// https://uber-go.github.io/fx/modules.html
	Module = fx.Module( //nolint:gochecknoglobals
		"oteltracingfx",
		tracingProviders,
		tracingInvokes,
	)

	// tracingProviders is a module for the tracing.
	tracingProviders = fx.Options(fx.Provide( //nolint:gochecknoglobals
		ProvideTracingConfig,
		NewOtelTracing,
		fx.Annotate(
			provideTracer,
			fx.ParamTags(`optional:"true"`),
			fx.As(new(AppTracer)),
			fx.As(new(trace.Tracer)),
		),
	))

	// tracingInvokes is a module for the tracing.
	tracingInvokes = fx.Options(
		fx.Invoke(registerHooks),
	)
)

// provideTracer provides a tracer.
func provideTracer(
	tracingOtel *TracingOpenTelemetry,
) AppTracer {
	return tracingOtel.appTracer
}

// registerHooks registers hooks for the tracing.
// we don't want to register any dependencies here, its func body should execute always even we don't request for that, so we should use `invoke`.
func registerHooks(
	lc fx.Lifecycle,
	logger logger.Logger,
	tracingOtel *TracingOpenTelemetry,
) {
	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			if err := tracingOtel.Shutdown(ctx); err != nil {
				logger.Errorf("error in shutting down trace provider: %v", err)
			} else {
				logger.Info("trace provider shutdown gracefully")
			}

			return nil
		},
	})
}
