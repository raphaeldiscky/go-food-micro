package infrastructure

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
	"go.uber.org/fx"

	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"

	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/configurations/rabbitmq"
)

// https://pmihaylov.com/shared-components-go-microservices/

// NewModule creates a new module.
func NewModule() fx.Option {
	return fx.Module(
		"infrastructurefx",
		// Modules
		core.Module,
		customEcho.Module,
		grpc.Module,
		mongodb.Module,
		redis.Module,
		rabbitmq.ModuleFunc(
			func(v *validator.Validate, l logger.Logger, tracer tracing.AppTracer) configurations.RabbitMQConfigurationBuilderFuc {
				return func(builder configurations.RabbitMQConfigurationBuilder) {
					rabbitmq2.ConfigProductsRabbitMQ(builder, l, v, tracer)
				}
			},
		),
		health.Module,
		tracing.Module,
		metrics.Module,

		// Other provides
		fx.Provide(validator.New),
	)
}
