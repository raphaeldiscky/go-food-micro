package infrastructure

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstroredb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/configurations"
	rabbitmq2 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/configurations/rabbitmq"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/params"

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
