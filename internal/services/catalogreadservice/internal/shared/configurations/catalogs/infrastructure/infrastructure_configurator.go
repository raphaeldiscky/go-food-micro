package infrastructure

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	loggingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/pipelines"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	metricspipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics/mediatr/pipelines"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	tracingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/mediatr/pipelines"

	"github.com/mehdihadeli/go-mediatr"
)

type InfrastructureConfigurator struct {
	contracts.Application
}

func NewInfrastructureConfigurator(
	app contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: app,
	}
}

func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
	ic.ResolveFunc(
		func(l logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics) error {
			err := mediatr.RegisterRequestPipelineBehaviors(
				loggingpipelines.NewMediatorLoggingPipeline(l),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(l),
				),
				metricspipelines.NewMediatorMetricsPipeline(
					metrics,
					metricspipelines.WithLogger(l),
				),
			)

			return err
		},
	)
}
