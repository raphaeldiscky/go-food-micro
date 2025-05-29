// Package products contains the products module.
package products

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/cqrs"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	"go.uber.org/fx"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/repositories"
	creatingproductv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/creatingproduct/v1"
	deletingproductv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/deletingproduct/v1"
	gettingproductbyidv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	gettingproductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproducts/v1"
	searchingproductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/searchingproduct/v1"
	updatingoroductsv1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/updatingproduct/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc"
)

// NewModule is a module that contains the products module.
func NewModule() fx.Option {
	return fx.Module(
		"productsfx",

		// Other provides
		fx.Provide(repositories.NewPostgresProductRepository),
		fx.Provide(grpc.NewProductGrpcService),

		fx.Provide(
			fx.Annotate(func(catalogsServer contracts.EchoHTTPServer) *echo.Group {
				var g *echo.Group
				catalogsServer.RouteBuilder().
					RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
						group := v1.Group("/products")
						g = group
					})

				return g
			}, fx.ResultTags(`name:"product-echo-group"`)),
		),

		// add cqrs handlers to DI
		fx.Provide(
			cqrs.AsHandler(
				creatingproductv1.NewCreateProductHandler,
				"product-handlers",
			),
			cqrs.AsHandler(
				gettingproductsv1.NewGetProductsHandler,
				"product-handlers",
			),
			cqrs.AsHandler(
				deletingproductv1.NewDeleteProductHandler,
				"product-handlers",
			),
			cqrs.AsHandler(
				gettingproductbyidv1.NewGetProductByIDHandler,
				"product-handlers",
			),
			cqrs.AsHandler(
				searchingproductsv1.NewSearchProductsHandler,
				"product-handlers",
			),
			cqrs.AsHandler(
				updatingoroductsv1.NewUpdateProductHandler,
				"product-handlers",
			),
		),

		// add endpoints to DI
		fx.Provide(
			route.AsRoute(
				creatingproductv1.NewCreteProductEndpoint,
				"product-routes",
			),
			route.AsRoute(
				updatingoroductsv1.NewUpdateProductEndpoint,
				"product-routes",
			),
			route.AsRoute(
				gettingproductsv1.NewGetProductsEndpoint,
				"product-routes",
			),
			route.AsRoute(
				searchingproductsv1.NewSearchProductsEndpoint,
				"product-routes",
			),
			route.AsRoute(
				gettingproductbyidv1.NewGetProductByIDEndpoint,
				"product-routes",
			),
			route.AsRoute(
				deletingproductv1.NewDeleteProductEndpoint,
				"product-routes",
			),
		),
	)
}
