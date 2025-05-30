// Package params contains the order projection parameters.
package params

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"go.uber.org/fx"
)

// OrderProjectionParams is the parameters for the order projection.
type OrderProjectionParams struct {
	fx.In

	Projections []projection.IProjection `group:"projections"`
}
