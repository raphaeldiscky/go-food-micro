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

// InfrastructureConfigurator is a struct that contains the infrastructure configurator.
type InfrastructureConfigurator struct {
	contracts.Application
}

// NewInfrastructureConfigurator is a constructor for the InfrastructureConfigurator.
func NewInfrastructureConfigurator(
	fxapp contracts.Application,
) *InfrastructureConfigurator {
	return &InfrastructureConfigurator{
		Application: fxapp,
	}
}

// ConfigInfrastructures is a method that configures the infrastructures.
func (ic *InfrastructureConfigurator) ConfigInfrastructures() {
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
