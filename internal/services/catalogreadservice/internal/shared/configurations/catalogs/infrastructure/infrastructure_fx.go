package infrastructure

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/configurations/rabbitmq"

	"github.com/go-playground/validator"
	"go.uber.org/fx"
)

// https://pmihaylov.com/shared-components-go-microservices/
var Module = fx.Module(
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
