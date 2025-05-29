// Package infrastructure contains the infrastructure fx.
package infrastructure

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstroredb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	"go.uber.org/fx"

	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"

	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/configurations/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/params"
)

// https://pmihaylov.com/shared-components-go-microservices/

// Module is the infrastructure fx module.
var Module = fx.Module(
	"infrastructurefx",
	// Modules
	core.Module,
	customEcho.Module,
	grpc.Module,
	mongodb.Module,
	elasticsearch.Module,
	eventstroredb.ModuleFunc(
		func(params params.OrderProjectionParams) eventstroredb.ProjectionBuilderFuc {
			return func(builder eventstroredb.ProjectionsBuilder) {
				builder.AddProjections(params.Projections)
			}
		},
	),
	rabbitmq.ModuleFunc(
		func() configurations.RabbitMQConfigurationBuilderFuc {
			return func(builder configurations.RabbitMQConfigurationBuilder) {
				rabbitmq2.ConfigOrdersRabbitMQ(builder)
			}
		},
	),
	health.Module,
	tracing.Module,
	metrics.Module,

	// Other provides
	fx.Provide(validator.New),
)
