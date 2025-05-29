// Package health provides a health check service.
package health

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// healthService is a struct that represents a health service.
type healthService struct {
	healthParams contracts.HealthParams
}

// NewHealthService is a function that creates a new health service.
func NewHealthService(
	healthParams contracts.HealthParams,
) contracts.HealthService {
	return &healthService{
		healthParams: healthParams,
	}
}

// CheckHealth is a function that checks the health.
func (service *healthService) CheckHealth(ctx context.Context) contracts.Check {
	checks := make(contracts.Check)

	for _, health := range service.healthParams.Healths {
		checks[health.GetHealthName()] = contracts.NewStatus(
			health.CheckHealth(ctx),
		)
	}

	return checks
}
