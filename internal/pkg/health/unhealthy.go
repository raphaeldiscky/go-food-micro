// Package health provides a unhealthy health check service.
package health

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// UnhealthyHealthService is a struct that represents a unhealthy health check service.
type UnhealthyHealthService struct{}

// NewUnhealthyHealthService is a function that creates a new unhealthy health check service.
func NewUnhealthyHealthService() UnhealthyHealthService {
	return UnhealthyHealthService{}
}

// CheckHealth is a function that checks the health.
func (service UnhealthyHealthService) CheckHealth(
	context.Context,
) contracts.Check {
	return contracts.Check{
		"postgres": contracts.Status{Status: contracts.StatusDown},
		"redis":    contracts.Status{Status: contracts.StatusDown},
	}
}
