// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"github.com/EventStore/EventStore-Client-Go/esdb"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
)

func NewEventStoreDB(cfg *config.EventStoreDbOptions) (*esdb.Client, error) {
	settings, err := esdb.ParseConnectionString(cfg.GrpcEndPoint())
	if err != nil {
		return nil, err
	}

	return esdb.NewClient(settings)
}
