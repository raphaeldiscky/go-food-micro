// Package contracts provides a health check contracts.
package contracts

import "context"

// Health is an interface that represents a health check.
type Health interface {
	CheckHealth(ctx context.Context) error
	GetHealthName() string
}

// HealthService is an interface that represents a health check service.
type HealthService interface {
	CheckHealth(ctx context.Context) Check
}
