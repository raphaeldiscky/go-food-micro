// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
)

// NewEventStoreDB creates a new event store db.
func NewEventStoreDB(cfg *config.EventStoreDbOptions) (*kdb.Client, error) {
	settings, err := kdb.ParseConnectionString(cfg.GrpcEndPoint())
	if err != nil {
		return nil, err
	}

	return kdb.NewClient(settings)
}
