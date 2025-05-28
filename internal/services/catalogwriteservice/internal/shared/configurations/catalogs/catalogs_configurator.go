// Package catalogs contains the catalogs service configurator.
package catalogs

import (
	"fmt"
	"net/http"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	"gorm.io/gorm"

	echo "github.com/labstack/echo/v4"
	echocontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	migrationcontracts "github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/configurations"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/configurations/catalogs/infrastructure"
)

// CatalogWriteServiceConfigurator is a struct that contains the catalogs service configurator.
type CatalogWriteServiceConfigurator struct {
	contracts.Application
	infrastructureConfigurator *infrastructure.CatalogWriteInfraConfigurator
	productsModuleConfigurator *configurations.ProductsModuleConfigurator
}

// NewCatalogWriteServiceConfigurator is a constructor for the CatalogWriteServiceConfigurator.
func NewCatalogWriteServiceConfigurator(
	app contracts.Application,
) *CatalogWriteServiceConfigurator {
	infraConfigurator := infrastructure.NewCatalogWriteInfraConfigurator(app)
	productModuleConfigurator := configurations.NewProductsModuleConfigurator(
		app,
	)

	return &CatalogWriteServiceConfigurator{
		Application:                app,
		infrastructureConfigurator: infraConfigurator,
		productsModuleConfigurator: productModuleConfigurator,
	}
}

// ConfigureCatalogs is a method that configures the catalogs.
func (ic *CatalogWriteServiceConfigurator) ConfigureCatalogs() error {
	ic.infrastructureConfigurator.CatalogWriteConfigInfra()

	// Catalogs configurations
	ic.ResolveFunc(
		func(db *gorm.DB, postgresMigrationRunner migrationcontracts.PostgresMigrationRunner) error {
			err := ic.migrateCatalogs(postgresMigrationRunner)
			if err != nil {
				return err
			}

			if ic.Environment() != environment.Test {
				err = ic.seedCatalogs(db)
				if err != nil {
					return err
				}
			}

			return nil
		},
	)

	// Modules
	// Product module
	err := ic.productsModuleConfigurator.ConfigureProductsModule()

	return err
}

// MapCatalogsEndpoints is a method that maps the catalogs endpoints.
func (ic *CatalogWriteServiceConfigurator) MapCatalogsEndpoints() error {
	// Shared
	ic.ResolveFunc(
		func(catalogsServer echocontracts.EchoHttpServer, options *config.AppOptions) error {
			catalogsServer.SetupDefaultMiddlewares()

			// config catalogs root endpoint
			catalogsServer.RouteBuilder().
				RegisterRoutes(func(e *echo.Echo) {
					e.GET("", func(ec echo.Context) error {
						return ec.String(
							http.StatusOK,
							fmt.Sprintf(
								"%s is running...",
								options.GetMicroserviceNameUpper(),
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
	err := ic.productsModuleConfigurator.MapProductsEndpoints()

	return err
}
