// Package health provides a health check module.
package health

import (
	"go.uber.org/fx"
)

// Module is a fx.Options that provides the health check module.
var Module = fx.Options(
	fx.Provide(
		NewHealthService,
		NewHealthCheckEndpoint,
	),
	// Invoke is a function that invokes the health check module.
	fx.Invoke(func(endpoint *HealthCheckEndpoint) {
		endpoint.RegisterEndpoints()
	}),
)
