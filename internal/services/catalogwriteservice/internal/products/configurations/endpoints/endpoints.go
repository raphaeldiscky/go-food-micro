// Package endpoints contains the products endpoints.
package endpoints

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
)

// RegisterEndpoints is a function that registers the endpoints.
func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	return nil
}
