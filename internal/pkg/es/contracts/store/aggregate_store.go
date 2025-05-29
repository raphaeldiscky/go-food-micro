package store

import (
	"context"

	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	appendResult "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/appendresult"
	readPosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/readposition"
	expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
)

// AggregateStore is responsible for loading and saving Aggregate.
type AggregateStore[T models.IHaveEventSourcedAggregate] interface {
	// StoreWithVersion store the new or update aggregate state with expected version
	StoreWithVersion(
		aggregate T,
		metadata metadata.Metadata,
		expectedVersion expectedStreamVersion.ExpectedStreamVersion,
		ctx context.Context) (*appendResult.AppendEventsResult, error)

	// Store the new or update aggregate state
	Store(
		aggregate T,
		metadata metadata.Metadata,
		ctx context.Context,
	) (*appendResult.AppendEventsResult, error)

	// Load loads the most recent version of an aggregate to provided  into params aggregate with an id and start read position.
	Load(ctx context.Context, aggregateID uuid.UUID) (T, error)

	// LoadWithReadPosition loads the most recent version of an aggregate to provided  into params aggregate with an id and read position.
	LoadWithReadPosition(
		ctx context.Context,
		aggregateID uuid.UUID,
		position readPosition.StreamReadPosition,
	) (T, error)

	// Exists check aggregate exists by AggregateId.
	Exists(ctx context.Context, aggregateID uuid.UUID) (bool, error)
}
