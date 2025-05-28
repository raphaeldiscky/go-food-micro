package catalogs

import (
	echo "github.com/labstack/echo/v4"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/docs"
)

func (ic *CatalogReadServiceConfigurator) configSwagger(routeBuilder *customEcho.RouteBuilder) {
	// https://github.com/swaggo/swag#how-to-use-it-with-gin
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Read-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Read-Service Api."

	routeBuilder.RegisterRoutes(func(e *echo.Echo) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	})
}
