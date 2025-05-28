package infrastructure

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/goose"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresmessaging"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	"go.uber.org/fx"

	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"

	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/rabbitmq"
)

// https://pmihaylov.com/shared-components-go-microservices/

// Module is a module that contains the infrastructure module.
var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	postgresgorm.Module,
	postgresmessaging.Module,
	goose.Module,
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigProductsRabbitMQ(builder)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
