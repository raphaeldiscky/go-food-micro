// Package configurations contains the orders module configurator.
package configurations

import (
	"github.com/go-playground/validator"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	contracts2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	grpcServer "github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	echocontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	googleGrpc "google.golang.org/grpc"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/configurations/mappings"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/configurations/mediatr"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc"
	ordersservice "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/grpc/genproto"
)

// OrdersModuleConfigurator is the orders module configurator.
type OrdersModuleConfigurator struct {
	contracts2.Application
}

// NewOrdersModuleConfigurator creates a new orders module configurator.
func NewOrdersModuleConfigurator(
	app contracts2.Application,
) *OrdersModuleConfigurator {
	return &OrdersModuleConfigurator{
		Application: app,
	}
}

// ConfigureOrdersModule configures the orders module.
func (c *OrdersModuleConfigurator) ConfigureOrdersModule() {
	c.ResolveFunc(
		func(logger logger.Logger,
			_ echocontracts.EchoHTTPServer,
			orderRepository repositories.OrderMongoRepository,
			orderAggregateStore store.AggregateStore[*aggregate.Order],
			tracer tracing.AppTracer,
		) error {
			// config Orders Mappings
			err := mappings.ConfigureOrdersMappings()
			if err != nil {
				return err
			}

			// config Orders Mediators
			err = mediatr.ConfigOrdersMediator(logger, orderRepository, orderAggregateStore, tracer)
			if err != nil {
				return err
			}

			return nil
		},
	)
}

// MapOrdersEndpoints maps the orders endpoints.
func (c *OrdersModuleConfigurator) MapOrdersEndpoints() {
	// config Orders Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"order-routes"`,
	)

	// config Orders Grpc Endpoints
	c.ResolveFunc(
		func(ordersGrpcServer grpcServer.GrpcServer, ordersMetrics *contracts.OrdersMetrics, logger logger.Logger, validator *validator.Validate) error {
			orderGrpcService := grpc.NewOrderGrpcService(logger, validator, ordersMetrics)
			ordersGrpcServer.GrpcServiceBuilder().RegisterRoutes(func(server *googleGrpc.Server) {
				ordersservice.RegisterOrdersServiceServer(server, orderGrpcService)
			})

			return nil
		},
	)
}
