// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"context"
	"fmt"
	"math"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/trace"

	linq "github.com/ahmetb/go-linq/v3"
	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	appendResult "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/appendresult"
	streamName "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamname"
	readPosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/readposition"
	truncatePosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/truncateposition"
	expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
	esErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
)

// eventStoreDbEventStore is a struct that represents a event store db event store.
// https://developers.eventstore.com/clients/grpc/reading-events.html#reading-from-a-stream
// https://developers.eventstore.com/clients/grpc/appending-events.html#append-your-first-event
type eventStoreDbEventStore struct {
	log        logger.Logger
	client     *kdb.Client
	serializer *EsdbSerializer
	tracer     trace.Tracer
}

// NewEventStoreDbEventStore creates a new event store db event store.
func NewEventStoreDbEventStore(
	log logger.Logger,
	client *kdb.Client,
	serializer *EsdbSerializer,
	tracer trace.Tracer,
) store.EventStore {
	return &eventStoreDbEventStore{
		log:        log,
		client:     client,
		serializer: serializer,
		tracer:     tracer,
	}
}

// StreamExists checks if a stream exists.
func (e *eventStoreDbEventStore) StreamExists(
	streamName streamName.StreamName,
	ctx context.Context,
) (bool, error) {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.StreamExists")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	stream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		kdb.ReadStreamOptions{
			Direction: kdb.Backwards,
			From:      kdb.End{},
		},
		1)
	if err != nil {
		return false, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewReadStreamError(err),
				"error in reading stream",
			),
		)
	}

	defer stream.Close()

	return stream != nil, nil
}

// AppendEvents appends events to a stream.
func (e *eventStoreDbEventStore) AppendEvents(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	events []*models.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.AppendEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	var eventsData []kdb.EventData
	linq.From(events).SelectT(func(s *models.StreamEvent) kdb.EventData {
		data, err := e.serializer.StreamEventToEventData(s)
		if err != nil {
			return *new(kdb.EventData)
		}

		return data
	}).ToSlice(&eventsData)

	var appendEventsResult *appendResult.AppendEventsResult

	res, err := e.client.AppendToStream(
		ctx,
		streamName.String(),
		kdb.AppendToStreamOptions{
			StreamState: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(
				expectedVersion,
			),
		},
		eventsData...)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewAppendToStreamError(err, streamName.String()),
				"error in appending to stream",
			),
		)
	}

	appendEventsResult = e.serializer.EsdbWriteResultToAppendEventResult(res)

	span.SetAttributes(
		attribute.Object("AppendEventsResult", appendEventsResult),
	)

	e.log.Infow(
		"events append to stream successfully",
		logger.Fields{
			"AppendEventsResult": appendEventsResult,
			"StreamId":           streamName.String(),
		},
	)

	return appendEventsResult, nil
}

// AppendNewEvents appends new events to a stream.
func (e *eventStoreDbEventStore) AppendNewEvents(
	streamName streamName.StreamName,
	events []*models.StreamEvent,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.AppendNewEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	appendEventsResult, err := e.AppendEvents(
		streamName,
		expectedStreamVersion.NoStream,
		events,
		ctx,
	)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewAppendToStreamError(err, streamName.String()),
				"error in appending to stream",
			),
		)
	}

	span.SetAttributes(attribute.Object("AppendNewEvents", appendEventsResult))

	e.log.Infow(
		"events append to stream successfully",
		logger.Fields{
			"AppendEventsResult": appendEventsResult,
			"StreamId":           streamName.String(),
		},
	)

	return appendEventsResult, nil
}

// ReadEvents reads events from a stream.
func (e *eventStoreDbEventStore) ReadEvents(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.ReadEvents")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	readStream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		kdb.ReadStreamOptions{
			Direction: kdb.Forwards,
			From: e.serializer.StreamReadPositionToStreamPosition(
				readPosition,
			),
			ResolveLinkTos: true,
		},
		count)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewReadStreamError(err),
				"error in reading stream",
			),
		)
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(
		readStream,
	)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in converting to resolved events",
			),
		)
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"error in converting to stream events",
			),
		)
	}

	return events, nil
}

// ReadEventsWithMaxCount reads events from a stream with a max count.
func (e *eventStoreDbEventStore) ReadEventsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(
		ctx,
		"eventStoreDbEventStore.ReadEventsWithMaxCount",
	)
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEvents(streamName, readPosition, uint64(math.MaxUint64), ctx)
}

// ReadEventsFromStart reads events from a stream from the start.
func (e *eventStoreDbEventStore) ReadEventsFromStart(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(
		ctx,
		"eventStoreDbEventStore.ReadEventsFromStart",
	)
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEvents(streamName, readPosition.Start, count, ctx)
}

// ReadEventsBackwards reads events from a stream backwards.
func (e *eventStoreDbEventStore) ReadEventsBackwards(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(
		ctx,
		"eventStoreDbEventStore.ReadEventsBackwards",
	)
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	readStream, err := e.client.ReadStream(
		ctx,
		streamName.String(),
		kdb.ReadStreamOptions{
			Direction: kdb.Backwards,
			From: e.serializer.StreamReadPositionToStreamPosition(
				readPosition,
			),
			ResolveLinkTos: true,
		},
		count)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewReadStreamError(err),
				"[eventStoreDbEventStore_ReadEventsBackwards:ReadStream] error in reading stream",
			),
		)
	}

	defer readStream.Close()

	resolvedEvents, err := e.serializer.EsdbReadStreamToResolvedEvents(
		readStream,
	)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[eventStoreDbEventStore_ReadEvents.EsdbReadStreamToResolvedEvents] error in converting to resolved events",
			),
		)
	}

	events, err := e.serializer.ResolvedEventsToStreamEvents(resolvedEvents)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIf(
				err,
				"[eventStoreDbEventStore_ReadEvents.ResolvedEventsToStreamEvents] error in converting to stream events",
			),
		)
	}

	return events, nil
}

// ReadEventsBackwardsWithMaxCount reads events from a stream backwards with a max count.
func (e *eventStoreDbEventStore) ReadEventsBackwardsWithMaxCount(
	streamName streamName.StreamName,
	readPosition readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(
		ctx,
		"eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount",
	)
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEventsBackwards(
		streamName,
		readPosition,
		uint64(math.MaxUint64),
		ctx,
	)
}

// ReadEventsBackwardsFromEnd reads events from a stream backwards from the end.
func (e *eventStoreDbEventStore) ReadEventsBackwardsFromEnd(
	streamName streamName.StreamName,
	count uint64,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	ctx, span := e.tracer.Start(
		ctx,
		"eventStoreDbEventStore.ReadEventsBackwardsWithMaxCount",
	)
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	return e.ReadEventsBackwards(streamName, readPosition.End, count, ctx)
}

func (e *eventStoreDbEventStore) TruncateStream(
	streamName streamName.StreamName,
	truncatePosition truncatePosition.StreamTruncatePosition,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.TruncateStream")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	streamMetadata := kdb.StreamMetadata{}
	streamMetadata.SetTruncateBefore(
		e.serializer.StreamTruncatePositionToInt64(truncatePosition),
	)
	writeResult, err := e.client.SetStreamMetadata(
		ctx,
		streamName.String(),
		kdb.AppendToStreamOptions{
			StreamState: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(
				expectedVersion,
			),
		},
		streamMetadata)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewTruncateStreamError(err, streamName.String()),
				"error in truncating stream",
			),
		)
	}

	span.SetAttributes(attribute.Object("WriteResult", writeResult))

	e.log.Infow(
		fmt.Sprintf(
			"stream with id %s truncated successfully",
			streamName.String(),
		),
		logger.Fields{
			"WriteResult": writeResult,
			"StreamId":    streamName.String(),
		},
	)

	return e.serializer.EsdbWriteResultToAppendEventResult(writeResult), nil
}

func (e *eventStoreDbEventStore) DeleteStream(
	streamName streamName.StreamName,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) error {
	ctx, span := e.tracer.Start(ctx, "eventStoreDbEventStore.DeleteStream")
	span.SetAttributes(attribute2.String("StreamName", streamName.String()))
	defer span.End()

	deleteResult, err := e.client.DeleteStream(
		ctx,
		streamName.String(),
		kdb.DeleteStreamOptions{
			StreamState: e.serializer.ExpectedStreamVersionToEsdbExpectedRevision(
				expectedVersion,
			),
		})
	if err != nil {
		return utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewDeleteStreamError(err, streamName.String()),
				"error in deleting stream",
			),
		)
	}

	span.SetAttributes(attribute.Object("DeleteResult", deleteResult))

	e.log.Infow(
		fmt.Sprintf(
			"stream with id %s deleted successfully",
			streamName.String(),
		),
		logger.Fields{
			"DeleteResult": deleteResult,
			"StreamId":     streamName.String(),
		},
	)

	return nil
}
