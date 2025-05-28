// Package infrastructure contains the infrastructure configurator.
package infrastructure

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"gorm.io/gorm"

	mediatr "github.com/mehdihadeli/go-mediatr"
	loggingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/pipelines"
	metricspipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics/mediatr/pipelines"
	tracingpipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/mediatr/pipelines"
	postgrespipelines "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/pipelines"
	validationpieline "github.com/raphaeldiscky/go-food-micro/internal/pkg/validation/pipeline"
)

// CatalogWriteInfraConfigurator is a struct that contains the infrastructure configurator.
type CatalogWriteInfraConfigurator struct {
	contracts.Application
}

// NewCatalogWriteInfraConfigurator is a constructor for the CatalogWriteInfraConfigurator.
func NewCatalogWriteInfraConfigurator(
	fxapp contracts.Application,
) *CatalogWriteInfraConfigurator {
	return &CatalogWriteInfraConfigurator{
		Application: fxapp,
	}
}

// CatalogWriteConfigInfra is a method that configures the infrastructures.
func (ic *CatalogWriteInfraConfigurator) CatalogWriteConfigInfra() {
	ic.ResolveFunc(
		func(l logger.Logger, tracer tracing.AppTracer, metrics metrics.AppMetrics, db *gorm.DB) error {
			err := mediatr.RegisterRequestPipelineBehaviors(
				loggingpipelines.NewMediatorLoggingPipeline(l),
				validationpieline.NewMediatorValidationPipeline(l),
				tracingpipelines.NewMediatorTracingPipeline(
					tracer,
					tracingpipelines.WithLogger(l),
				),
				metricspipelines.NewMediatorMetricsPipeline(
					metrics,
					metricspipelines.WithLogger(l),
				),
				postgrespipelines.NewMediatorTransactionPipeline(l, db),
			)

			return err
		},
	)
}
