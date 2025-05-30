// Package contracts provides a eventstoredb container contracts.
package contracts

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
)

// EventstoreDBContainerOptions represents a eventstoredb container options.
type EventstoreDBContainerOptions struct {
	Host    string
	Ports   []string
	TcpPort int
	// HTTP is the primary protocol for EventStoreDB. It is used in gRPC communication and HTTP APIs (management, gossip and diagnostics).
	HttpPort  int
	ImageName string
	Name      string
	Tag       string
}

// EventstoreDBContainer is a interface that represents a eventstoredb container.
type EventstoreDBContainer interface {
	PopulateContainerOptions(
		ctx context.Context,
		t *testing.T,
		options ...*EventstoreDBContainerOptions,
	) (*config.EventStoreDbOptions, error)

	Cleanup(ctx context.Context) error
}
