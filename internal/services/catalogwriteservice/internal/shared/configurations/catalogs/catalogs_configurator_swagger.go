package catalogs

import (
	echo "github.com/labstack/echo/v4"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/docs"
)

func (ic *CatalogsServiceConfigurator) configSwagger(routeBuilder *customEcho.RouteBuilder) {
	// https://github.com/swaggo/swag#how-to-use-it-with-gin
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Catalogs Write-Service Api"
	docs.SwaggerInfo.Description = "Catalogs Write-Service Api."

	routeBuilder.RegisterRoutes(func(e *echo.Echo) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	})
}
