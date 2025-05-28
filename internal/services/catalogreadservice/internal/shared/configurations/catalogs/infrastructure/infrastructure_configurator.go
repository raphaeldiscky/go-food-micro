// Package infrastructure contains the infrastructure configurator.
package infrastructure

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	mediatr "github.com/mehdihadeli/go-mediatr"
	loggingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/pipelines"
	metricspipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics/mediatr/pipelines"
	tracingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/mediatr/pipelines"
)

// CatalogReadInfraConfigurator is a struct that contains the infrastructure configurator.
type CatalogReadInfraConfigurator struct {
	contracts.Application
}

// NewCatalogReadInfraConfigurator creates a new CatalogReadInfraConfigurator.
func NewCatalogReadInfraConfigurator(
	app contracts.Application,
) *CatalogReadInfraConfigurator {
	return &CatalogReadInfraConfigurator{
		Application: app,
	}
}

// CatalogReadConfigInfra is a method that configures the infrastructures.
func (ic *CatalogReadInfraConfigurator) CatalogReadConfigInfra() {
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
