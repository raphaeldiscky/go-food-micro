package orders

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstroredb"
	"go.uber.org/fx"

	echo "github.com/labstack/echo/v4"
	echocontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/data/repositories"
	createOrderV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/creatingorder/v1/endpoints"
	GetOrderByIDV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorderbyid/v1/endpoints"
	getOrdersV1 "github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/features/gettingorders/v1/endpoints"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/models/orders/aggregate"
	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/projections"
)

var Module = fx.Module(
	"ordersfx",

	// Other provides
	fx.Provide(fx.Annotate(repositories.NewMongoOrderReadRepository)),
	fx.Provide(repositories.NewElasticOrderReadRepository),

	fx.Provide(eventstroredb.NewEventStoreAggregateStore[*aggregate.Order]),
	fx.Provide(fx.Annotate(func(catalogsServer echocontracts.EchoHttpServer) *echo.Group {
		var g *echo.Group
		catalogsServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
			group := v1.Group("/orders")
			g = group
		})

		return g
	}, fx.ResultTags(`name:"order-echo-group"`))),

	fx.Provide(
		route.AsRoute(createOrderV1.NewCreteOrderEndpoint, "order-routes"),
		route.AsRoute(GetOrderByIDV1.NewGetOrderByIDEndpoint, "order-routes"),
		route.AsRoute(getOrdersV1.NewGetOrdersEndpoint, "order-routes"),
	),

	fx.Provide(
		es.AsProjection(projections.NewElasticOrderProjection),
		es.AsProjection(projections.NewMongoOrderProjection),
	),
)
