package health

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	contracts2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
)

type HealthCheckEndpoint struct {
	service    contracts2.HealthService
	echoServer contracts.EchoHttpServer
}

func NewHealthCheckEndpoint(
	service contracts2.HealthService,
	server contracts.EchoHttpServer,
) *HealthCheckEndpoint {
	return &HealthCheckEndpoint{service: service, echoServer: server}
}

func (s *HealthCheckEndpoint) RegisterEndpoints() {
	s.echoServer.GetEchoInstance().GET("health", s.checkHealth)
}

func (s *HealthCheckEndpoint) checkHealth(c echo.Context) error {
	check := s.service.CheckHealth(c.Request().Context())
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}

	return c.JSON(http.StatusOK, check)
}
