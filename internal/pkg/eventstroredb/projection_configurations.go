package eventstroredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
