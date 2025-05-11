package params

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"

	"go.uber.org/fx"
)

type OrderProjectionParams struct {
	fx.In

	Projections []projection.IProjection `group:"projections"`
}
