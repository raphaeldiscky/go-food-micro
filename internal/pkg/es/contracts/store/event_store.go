// Package store provides the event store.
package store

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	appendResult "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/appendresult"
	streamName "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamname"
	readPosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/readposition"
	truncatePosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/truncateposition"
	expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
)

// EventStore is a interface that represents the event store.
type EventStore interface {
	// StreamExists Check if specific stream exists in the store
	StreamExists(streamName streamName.StreamName, ctx context.Context) (bool, error)

	// ReadEventsFromStart Read events for existing stream and start position with specified events count.
	ReadEventsFromStart(
		streamName streamName.StreamName,
		count uint64,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// ReadEvents Read events for a specific stream and position in forward mode with specified events count.
	ReadEvents(
		streamName streamName.StreamName,
		readPosition readPosition.StreamReadPosition,
		count uint64,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// ReadEventsWithMaxCount Read events for a specific stream and position in forward mode with max count.
	ReadEventsWithMaxCount(
		streamName streamName.StreamName,
		readPosition readPosition.StreamReadPosition,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// ReadEventsBackwards Read events for a specific stream and position in backwards mode with specified events count.
	ReadEventsBackwards(
		streamName streamName.StreamName,
		readPosition readPosition.StreamReadPosition,
		count uint64,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// ReadEventsBackwardsFromEnd Read events for a specific stream and end position in backwards mode with specified events count.
	ReadEventsBackwardsFromEnd(
		streamName streamName.StreamName,
		count uint64,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// ReadEventsBackwardsWithMaxCount Read events for a specific stream and position in backwards mode with max events count.
	ReadEventsBackwardsWithMaxCount(
		stream streamName.StreamName,
		readPosition readPosition.StreamReadPosition,
		ctx context.Context,
	) ([]*models.StreamEvent, error)

	// AppendEvents Append events to aggregate with an existing or none existing stream.
	AppendEvents(
		streamName streamName.StreamName,
		expectedVersion expectedStreamVersion.ExpectedStreamVersion,
		events []*models.StreamEvent,
		ctx context.Context,
	) (*appendResult.AppendEventsResult, error)

	// AppendNewEvents Append events to aggregate with none existing stream.
	AppendNewEvents(
		streamName streamName.StreamName,
		events []*models.StreamEvent,
		ctx context.Context,
	) (*appendResult.AppendEventsResult, error)

	// TruncateStream Truncate a stream at a given position
	TruncateStream(
		streamName streamName.StreamName,
		truncatePosition truncatePosition.StreamTruncatePosition,
		expectedVersion expectedStreamVersion.ExpectedStreamVersion,
		ctx context.Context,
	) (*appendResult.AppendEventsResult, error)

	// DeleteStream Delete a stream
	DeleteStream(
		streamName streamName.StreamName,
		expectedVersion expectedStreamVersion.ExpectedStreamVersion,
		ctx context.Context,
	) error
}
