// Package configurations contains the products module configurator.
package configurations

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing"

	logger2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/configurations/mappings"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/configurations/mediator"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
)

// ProductsModuleConfigurator is a struct that contains the products module configurator.
type ProductsModuleConfigurator struct {
	contracts.Application
}

// NewProductsModuleConfigurator is a constructor for the ProductsModuleConfigurator.
func NewProductsModuleConfigurator(
	app contracts.Application,
) *ProductsModuleConfigurator {
	return &ProductsModuleConfigurator{
		Application: app,
	}
}

// ConfigureProductsModule is a method that configures the products module.
func (c *ProductsModuleConfigurator) ConfigureProductsModule() {
	c.ResolveFunc(
		func(logger logger2.Logger, mongoRepository data.ProductRepository, cacheRepository data.ProductCacheRepository, tracer tracing.AppTracer) error {
			// config Products Mediators
			err := mediator.ConfigProductsMediator(
				logger,
				mongoRepository,
				cacheRepository,
				tracer,
			)
			if err != nil {
				return err
			}

			// config Products Mappings
			err = mappings.ConfigureProductsMappings()
			if err != nil {
				return err
			}

			return nil
		},
	)
}

// MapProductsEndpoints is a method that maps the products endpoints.
func (c *ProductsModuleConfigurator) MapProductsEndpoints() {
	// config Products Http Endpoints
	c.ResolveFuncWithParamTag(func(endpoints []route.Endpoint) {
		for _, endpoint := range endpoints {
			endpoint.MapEndpoint()
		}
	}, `group:"product-routes"`,
	)
}
