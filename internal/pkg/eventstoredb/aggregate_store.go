// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"context"
	"fmt"
	"reflect"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/trace"

	linq "github.com/ahmetb/go-linq/v3"
	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	uuid "github.com/satori/go.uuid"
	attribute2 "go.opentelemetry.io/otel/attribute"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/store"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	appendResult "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/appendresult"
	streamName "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamname"
	readPosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/readposition"
	expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
	esErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/attribute"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/tracing/utils"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// esdbAggregateStore is a struct that represents a event store aggregate store.
type esdbAggregateStore[T models.IHaveEventSourcedAggregate] struct {
	log        logger.Logger
	eventStore store.EventStore
	serializer *EsdbSerializer
	tracer     trace.Tracer
}

// NewEventStoreAggregateStore creates a new event store aggregate store.
func NewEventStoreAggregateStore[T models.IHaveEventSourcedAggregate](
	log logger.Logger,
	eventStore store.EventStore,
	serializer *EsdbSerializer,
	tracer trace.Tracer,
) store.AggregateStore[T] {
	return &esdbAggregateStore[T]{
		log:        log,
		eventStore: eventStore,
		serializer: serializer,
		tracer:     tracer,
	}
}

// StoreWithVersion stores an aggregate with a version.
func (a *esdbAggregateStore[T]) StoreWithVersion(
	aggregate T,
	metadata metadata.Metadata,
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := a.tracer.Start(ctx, "esdbAggregateStore.StoreWithVersion")
	span.SetAttributes(
		attribute2.String("AggregateID", aggregate.ID().String()),
	)
	defer span.End()

	if len(aggregate.UncommittedEvents()) == 0 {
		a.log.Infow(
			fmt.Sprintf(
				"[esdbAggregateStore.StoreWithVersion] No events to store for aggregateID %s",
				aggregate.ID(),
			),
			logger.Fields{"AggregateID": aggregate.ID()},
		)

		return appendResult.NoOp, nil
	}

	streamId := streamName.For[T](aggregate)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	var streamEvents []*models.StreamEvent

	linq.From(aggregate.UncommittedEvents()).
		SelectIndexedT(func(i int, domainEvent domain.IDomainEvent) *models.StreamEvent {
			var inInterface map[string]interface{}
			err := a.serializer.eventSerializer.Serializer().
				DecodeWithMapStructure(domainEvent, &inInterface)
			if err != nil {
				return nil
			}

			return a.serializer.DomainEventToStreamEvent(
				domainEvent,
				metadata,
				int64(i)+aggregate.OriginalVersion(),
			)
		}).
		ToSlice(&streamEvents)

	streamAppendResult, err := a.eventStore.AppendEvents(
		streamId,
		expectedVersion,
		streamEvents,
		ctx,
	)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIff(
				err,
				"[esdbAggregateStore_StoreWithVersion:AppendEvents] error in storing aggregate with id {%d}",
				aggregate.ID(),
			),
		)
	}

	aggregate.MarkUncommittedEventAsCommitted()

	span.SetAttributes(attribute.Object("Aggregate", aggregate))

	a.log.Infow(
		fmt.Sprintf(
			"[esdbAggregateStore.StoreWithVersion] aggregate with id %d stored successfully",
			aggregate.ID(),
		),
		logger.Fields{"Aggregate": aggregate, "StreamId": streamId},
	)

	return streamAppendResult, nil
}

// Store stores an aggregate.
func (a *esdbAggregateStore[T]) Store(
	aggregate T,
	metadata metadata.Metadata,
	ctx context.Context,
) (*appendResult.AppendEventsResult, error) {
	ctx, span := a.tracer.Start(ctx, "esdbAggregateStore.Store")
	defer span.End()

	expectedVersion := expectedStreamVersion.FromInt64(
		aggregate.OriginalVersion(),
	)

	streamAppendResult, err := a.StoreWithVersion(
		aggregate,
		metadata,
		expectedVersion,
		ctx,
	)
	if err != nil {
		return nil, utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIff(
				err,
				"[esdbAggregateStore_Store:StoreWithVersion] failed to store aggregate with id{%v}",
				aggregate.ID(),
			),
		)
	}

	return streamAppendResult, nil
}

// Load loads an aggregate.
func (a *esdbAggregateStore[T]) Load(
	ctx context.Context,
	aggregateID uuid.UUID,
) (T, error) {
	ctx, span := a.tracer.Start(ctx, "esdbAggregateStore.Load")
	defer span.End()

	position := readPosition.Start

	return a.LoadWithReadPosition(ctx, aggregateID, position)
}

// LoadWithReadPosition loads an aggregate with a read position.
func (a *esdbAggregateStore[T]) LoadWithReadPosition(
	ctx context.Context,
	aggregateID uuid.UUID,
	position readPosition.StreamReadPosition,
) (T, error) {
	ctx, span := a.tracer.Start(ctx, "esdbAggregateStore.LoadWithReadPosition")
	span.SetAttributes(attribute2.String("AggregateID", aggregateID.String()))
	defer span.End()

	var typeNameType T
	aggregateInstance := typeMapper.InstancePointerByTypeName(
		typeMapper.GetFullTypeName(typeNameType),
	)
	aggregate, ok := aggregateInstance.(T)
	if !ok {
		return *new(T), errors.New(
			fmt.Sprintf(
				"[esdbAggregateStore_LoadWithReadPosition] aggregate is not a %s",
				typeMapper.GetFullTypeName(typeNameType),
			),
		)
	}

	method := reflect.ValueOf(aggregate).MethodByName("NewEmptyAggregate")
	if !method.IsValid() {
		return *new(T), utils.TraceErrStatusFromSpan(
			span,
			errors.New(
				"[esdbAggregateStore_LoadWithReadPosition:MethodByName] aggregate does not have a `NewEmptyAggregate` method",
			),
		)
	}

	method.Call([]reflect.Value{})

	streamId := streamName.ForID[T](aggregateID)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	streamEvents, err := a.getStreamEvents(streamId, position, ctx)
	var kdbErr *kdb.Error
	if (errors.As(err, &kdbErr) && kdbErr.Code() == kdb.ErrorCodeResourceNotFound) ||
		len(streamEvents) == 0 {
		return *new(T), utils.TraceErrStatusFromSpan(
			span,
			errors.WithMessage(
				esErrors.NewAggregateNotFoundError(err, aggregateID),
				"[esdbAggregateStore.LoadWithReadPosition] error in loading aggregate",
			),
		)
	}

	if err != nil {
		return *new(T), utils.TraceErrStatusFromSpan(
			span,
			errors.WrapIff(
				err,
				"[esdbAggregateStore.LoadWithReadPosition:MethodByName] error in loading aggregate {%s}",
				aggregateID.String(),
			),
		)
	}

	var meta metadata.Metadata
	var domainEvents []domain.IDomainEvent

	linq.From(streamEvents).
		Distinct().
		SelectT(func(streamEvent *models.StreamEvent) domain.IDomainEvent {
			meta = streamEvent.Metadata

			return streamEvent.Event
		}).
		ToSlice(&domainEvents)

	err = aggregate.LoadFromHistory(domainEvents, meta)
	if err != nil {
		return *new(T), utils.TraceStatusFromSpan(span, err)
	}

	a.log.Infow(
		fmt.Sprintf("Loaded aggregate with streamId {%s} and aggregateID {%s}",
			streamId.String(),
			aggregateID.String()),
		logger.Fields{"Aggregate": aggregate, "StreamId": streamId.String()},
	)

	span.SetAttributes(attribute.Object("Aggregate", aggregate))

	return aggregate, nil
}

// Exists checks if an aggregate exists.
func (a *esdbAggregateStore[T]) Exists(
	ctx context.Context,
	aggregateID uuid.UUID,
) (bool, error) {
	ctx, span := a.tracer.Start(ctx, "esdbAggregateStore.Exists")
	span.SetAttributes(attribute2.String("AggregateID", aggregateID.String()))
	defer span.End()

	streamId := streamName.ForID[T](aggregateID)
	span.SetAttributes(attribute2.String("StreamId", streamId.String()))

	return a.eventStore.StreamExists(streamId, ctx)
}

// getStreamEvents gets stream events.
func (a *esdbAggregateStore[T]) getStreamEvents(
	streamId streamName.StreamName,
	position readPosition.StreamReadPosition,
	ctx context.Context,
) ([]*models.StreamEvent, error) {
	pageSize := 500
	var streamEvents []*models.StreamEvent

	for {
		events, err := a.eventStore.ReadEvents(
			streamId,
			position,
			//nolint:gosec // G115: integer overflow conversion int -> uint64
			uint64(pageSize),
			ctx,
		)
		if err != nil {
			return nil, errors.WrapIff(
				err,
				"[esdbAggregateStore_getStreamEvents:ReadEvents] failed to read events",
			)
		}
		streamEvents = append(streamEvents, events...)
		if len(events) < pageSize {
			break
		}
		position = readPosition.FromInt64(int64(len(events)) + position.Value())
	}

	return streamEvents, nil
}
