// Package products contains the products module.
package products

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	"go.uber.org/fx"

	echo "github.com/labstack/echo/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/data/repositories"
	getProductByIdV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/getproductbyid/v1/endpoints"
	getProductsV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/gettingproducts/v1/endpoints"
	searchProductV1 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/searchingproducts/v1/endpoints"
)

// NewModule is a module that contains the products module.
func NewModule() fx.Option {
	return fx.Module(
		"productsfx",

		// Other provides
		fx.Provide(repositories.NewRedisProductRepository),
		fx.Provide(repositories.NewMongoProductRepository),

		fx.Provide(fx.Annotate(func(catalogsServer contracts.EchoHTTPServer) *echo.Group {
			var g *echo.Group
			catalogsServer.RouteBuilder().RegisterGroupFunc("/api/v1", func(v1 *echo.Group) {
				group := v1.Group("/products")
				g = group
			})

			return g
		}, fx.ResultTags(`name:"product-echo-group"`))),

		fx.Provide(
			route.AsRoute(getProductsV1.NewGetProductsEndpoint, "product-routes"),
			route.AsRoute(searchProductV1.NewSearchProductsEndpoint, "product-routes"),
			route.AsRoute(getProductByIdV1.NewGetProductByIDEndpoint, "product-routes"),
		),
	)
}
