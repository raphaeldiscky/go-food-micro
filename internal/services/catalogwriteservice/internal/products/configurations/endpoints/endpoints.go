package endpoints

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/web/route"
)

func RegisterEndpoints(endpoints []route.Endpoint) error {
	for _, endpoint := range endpoints {
		endpoint.MapEndpoint()
	}

	return nil
}
