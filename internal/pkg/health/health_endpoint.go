// Package health provides a health check endpoint.
package health

import (
	"net/http"

	echo "github.com/labstack/echo/v4"

	contracts2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
)

// HealthCheckEndpoint is a struct that represents a health check endpoint.
type HealthCheckEndpoint struct {
	service    contracts2.HealthService
	echoServer contracts.EchoHTTPServer
}

// NewHealthCheckEndpoint is a function that creates a new health check endpoint.
func NewHealthCheckEndpoint(
	service contracts2.HealthService,
	server contracts.EchoHTTPServer,
) *HealthCheckEndpoint {
	return &HealthCheckEndpoint{service: service, echoServer: server}
}

// RegisterEndpoints is a function that registers the endpoints.
func (s *HealthCheckEndpoint) RegisterEndpoints() {
	s.echoServer.GetEchoInstance().GET("health", s.checkHealth)
}

// checkHealth is a function that checks the health.
func (s *HealthCheckEndpoint) checkHealth(c echo.Context) error {
	check := s.service.CheckHealth(c.Request().Context())
	if !check.AllUp() {
		return c.JSON(http.StatusServiceUnavailable, check)
	}

	return c.JSON(http.StatusOK, check)
}
