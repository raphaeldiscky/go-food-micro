// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

type ProjectionsBuilder interface {
	AddProjection(projection projection.IProjection) ProjectionsBuilder
	AddProjections(projections []projection.IProjection) ProjectionsBuilder
	Build() *ProjectionsConfigurations
}
