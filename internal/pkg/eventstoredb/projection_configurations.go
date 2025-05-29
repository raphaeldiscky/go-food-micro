// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

type ProjectionsConfigurations struct {
	Projections []projection.IProjection
}
