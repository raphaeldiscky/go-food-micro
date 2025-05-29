// Package configurations contains the products module configurator.
package configurations

import (
	fxcontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	grpcServer "github.com/raphaeldiscky/go-food-micro/internal/pkg/grpc"
	googleGrpc "google.golang.org/grpc"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/endpoints"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/mappings"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations/mediator"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc"
	productsservice "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/grpc/genproto"
)

// ProductsModuleConfigurator is a struct that contains the products module configurator.
type ProductsModuleConfigurator struct {
	fxcontracts.Application
}

// NewProductsModuleConfigurator is a constructor for the ProductsModuleConfigurator.
func NewProductsModuleConfigurator(
	fxapp fxcontracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: fxapp,
	}
}

// ConfigureProductsModule is a method that configures the products module.
func (c *ProductsModuleConfigurator) ConfigureProductsModule() error {
	// config products mappings
	err := mappings.ConfigureProductsMappings()
	if err != nil {
		return err
	}

	// register products request handler on mediator
	c.ResolveFuncWithParamTag(
		mediator.RegisterMediatorHandlers,
		`group:"product-handlers"`,
	)

	return nil
}

// MapProductsEndpoints is a method that maps the products endpoints.
func (c *ProductsModuleConfigurator) MapProductsEndpoints() error {
	// config endpoints
	c.ResolveFuncWithParamTag(
		endpoints.RegisterEndpoints,
		`group:"product-routes"`,
	)

	// config Products Grpc Endpoints
	c.ResolveFunc(
		func(catalogsGrpcServer grpcServer.GrpcServer, grpcService *grpc.ProductGrpcServiceServer) error {
			catalogsGrpcServer.GrpcServiceBuilder().
				RegisterRoutes(func(server *googleGrpc.Server) {
					productsservice.RegisterProductsServiceServer(
						server,
						grpcService,
					)
				})

			return nil
		},
	)

	return nil
}
