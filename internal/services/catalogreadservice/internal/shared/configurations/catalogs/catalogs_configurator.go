// Package catalogs contains the catalogs service configurator.
package catalogs

import (
	"fmt"
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"

	echo "github.com/labstack/echo/v4"
	echocontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs/infrastructure"
)

// CatalogReadServiceConfigurator is a struct that contains the catalogs service configurator.
type CatalogReadServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.CatalogReadInfraConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

// NewCatalogReadServiceConfigurator is a constructor for the CatalogReadServiceConfigurator.
func NewCatalogReadServiceConfigurator(app contracts.Application) *CatalogReadServiceConfigurator {
	infraConfigurator := infrastructure.NewCatalogReadInfraConfigurator(app)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(app)

	return &CatalogReadServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

// ConfigureCatalogs is a method that configures the catalogs.
func (ic *CatalogReadServiceConfigurator) ConfigureCatalogs() {
	ic.infrastructureConfigurator.CatalogReadConfigInfra()
	ic.productsModuleConfigurator.ConfigureProductsModule()
}

// MapCatalogsEndpoints is a method that maps the catalogs endpoints.
func (ic *CatalogReadServiceConfigurator) MapCatalogsEndpoints() {
	// Shared
	ic.ResolveFunc(
		func(catalogsServer echocontracts.EchoHTTPServer, cfg *config.Config) error {
			catalogsServer.SetupDefaultMiddlewares()

			// config catalogs root endpoint
			catalogsServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								cfg.AppOptions.GetMicroserviceNameUpper(),
							),
						)
					})
				})

			// config catalogs swagger
			ic.configSwagger(catalogsServer.RouteBuilder())

			return nil
		},
	)

	// Modules
	// Products CatalogsServiceModule endpoints
	ic.productsModuleConfigurator.MapProductsEndpoints()
}
