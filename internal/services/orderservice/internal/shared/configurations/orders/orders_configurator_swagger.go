package orders

import (
	echo "github.com/labstack/echo/v4"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	echoSwagger "github.com/swaggo/echo-swagger"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/docs"
)

func (ic *OrdersServiceConfigurator) configSwagger(routeBuilder *customEcho.RouteBuilder) {
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Title = "Orders Service Api"
	docs.SwaggerInfo.Description = "Orders Service Api."

	routeBuilder.RegisterRoutes(func(e *echo.Echo) {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	})
}
