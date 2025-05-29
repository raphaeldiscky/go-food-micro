// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

// ProjectionsConfigurations is a struct that represents a projections configurations.
type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
