// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"io"
	"strings"

	"emperror.dev/errors"

	linq "github.com/ahmetb/go-linq/v3"
	googleuuid "github.com/google/uuid"
	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
	appendResult "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/appendresult"
	readPosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/readposition"
	truncatePosition "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamposition/truncateposition"
	expectedStreamVersion "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
	esErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/errors"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// EsdbSerializer is a struct that represents a event store db serializer.
type EsdbSerializer struct {
	metadataSerializer serializer.MetadataSerializer
	eventSerializer    serializer.EventSerializer
}

// NewEsdbSerializer creates a new event store db serializer.
func NewEsdbSerializer(
	metadataSerializer serializer.MetadataSerializer,
	eventSerializer serializer.EventSerializer,
) *EsdbSerializer {
	return &EsdbSerializer{
		metadataSerializer: metadataSerializer,
		eventSerializer:    eventSerializer,
	}
}

// StreamEventToEventData converts a stream event to a event data.
func (e *EsdbSerializer) StreamEventToEventData(
	streamEvent *models.StreamEvent,
) (kdb.EventData, error) {
	eventSerializationResult, err := e.eventSerializer.Serialize(streamEvent.Event)
	if err != nil {
		return *new(kdb.EventData), err
	}

	metadataSerializationResult, err := e.metadataSerializer.Serialize(streamEvent.Metadata)
	if err != nil {
		return *new(kdb.EventData), err
	}

	var contentType kdb.ContentType

	switch eventSerializationResult.ContentType {
	case "application/json":
		contentType = kdb.ContentTypeJson
	default:
		contentType = kdb.ContentTypeBinary
	}

	id, err := uuid.FromString(streamEvent.EventID.String())
	if err != nil {
		return *new(kdb.EventData), err
	}

	googleID, err := googleuuid.Parse(id.String())
	if err != nil {
		return *new(kdb.EventData), err
	}

	return kdb.EventData{
		EventID:     googleID,
		EventType:   typeMapper.GetTypeName(streamEvent.Event),
		Data:        eventSerializationResult.Data,
		Metadata:    metadataSerializationResult,
		ContentType: contentType,
	}, nil
}

// ExpectedStreamVersionToEsdbExpectedRevision converts a expected stream version to a event store db expected revision.
func (e *EsdbSerializer) ExpectedStreamVersionToEsdbExpectedRevision(
	expectedVersion expectedStreamVersion.ExpectedStreamVersion,
) kdb.StreamState {
	if expectedVersion.IsNoStream() {
		return kdb.NoStream{}
	}
	if expectedVersion.IsAny() {
		return kdb.Any{}
	}
	if expectedVersion.IsStreamExists() {
		return kdb.StreamExists{}
	}

	//nolint:gosec // G115: integer overflow conversion int -> uint64
	return kdb.StreamRevision{Value: uint64(expectedVersion.Value())}
}

// StreamReadPositionToStreamPosition converts a stream read position to a stream position.
func (e *EsdbSerializer) StreamReadPositionToStreamPosition(
	readPosition readPosition.StreamReadPosition,
) kdb.StreamPosition {
	if readPosition.IsEnd() {
		return kdb.End{}
	}
	if readPosition.IsStart() {
		return kdb.Start{}
	}

	return kdb.Revision(1)
}

// StreamTruncatePositionToInt64 converts a stream truncate position to a int64.
func (e *EsdbSerializer) StreamTruncatePositionToInt64(
	truncatePosition truncatePosition.StreamTruncatePosition,
) uint64 {
	//nolint:gosec // G115: integer overflow conversion int -> uint64
	return uint64(truncatePosition.Value())
}

// EsdbReadStreamToResolvedEvents converts a event store db read stream to a resolved events.
func (e *EsdbSerializer) EsdbReadStreamToResolvedEvents(
	stream *kdb.ReadStream,
) ([]*kdb.ResolvedEvent, error) {
	var events []*kdb.ResolvedEvent

	for {
		event, err := stream.Recv()
		var kdbErr *kdb.Error
		if errors.As(err, &kdbErr) && kdbErr.Code() == kdb.ErrorCodeResourceNotFound {
			return nil, esErrors.NewStreamNotFoundError(err, event.Event.StreamID)
		}
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, esErrors.NewReadStreamError(err)
		}

		events = append(events, event)
	}

	return events, nil
}

// EsdbPositionToStreamReadPosition converts a event store db position to a stream read position.
func (e *EsdbSerializer) EsdbPositionToStreamReadPosition(
	position kdb.Position,
) readPosition.StreamReadPosition {
	//nolint:gosec // G115: integer overflow conversion int -> int64
	return readPosition.FromInt64(int64(position.Commit))
}

// ResolvedEventToStreamEvent converts a resolved event to a stream event.
func (e *EsdbSerializer) ResolvedEventToStreamEvent(
	resolveEvent *kdb.ResolvedEvent,
) (*models.StreamEvent, error) {
	deserializedEvent, err := e.eventSerializer.Deserialize(
		resolveEvent.Event.Data,
		resolveEvent.Event.EventType,
		resolveEvent.Event.ContentType,
	)
	if err != nil {
		return nil, err
	}

	deserializedMeta, err := e.metadataSerializer.Deserialize(resolveEvent.Event.UserMetadata)
	if err != nil {
		return nil, err
	}

	id, err := uuid.FromString(resolveEvent.Event.EventID.String())
	if err != nil {
		return nil, err
	}

	return &models.StreamEvent{
		EventID:  id,
		Event:    deserializedEvent,
		Metadata: deserializedMeta,
		//nolint:gosec // G115: integer overflow conversion int -> int64
		Version:  int64(resolveEvent.Event.EventNumber),
		Position: e.EsdbPositionToStreamReadPosition(resolveEvent.OriginalEvent().Position).Value(),
	}, nil
}

// ResolvedEventsToStreamEvents converts a resolved events to a stream events.
func (e *EsdbSerializer) ResolvedEventsToStreamEvents(
	resolveEvents []*kdb.ResolvedEvent,
) ([]*models.StreamEvent, error) {
	var streamEvents []*models.StreamEvent

	linq.From(resolveEvents).WhereT(func(item *kdb.ResolvedEvent) bool {
		return !strings.HasPrefix(item.Event.EventType, "$")
	}).SelectT(func(item *kdb.ResolvedEvent) *models.StreamEvent {
		event, err := e.ResolvedEventToStreamEvent(item)
		if err != nil {
			return nil
		}

		return event
	}).ToSlice(&streamEvents)

	return streamEvents, nil
}

// EsdbWriteResultToAppendEventResult converts a event store db write result to a append event result.
func (e *EsdbSerializer) EsdbWriteResultToAppendEventResult(
	writeResult *kdb.WriteResult,
) *appendResult.AppendEventsResult {
	return appendResult.From(writeResult.CommitPosition, writeResult.NextExpectedVersion)
}

// Serialize serializes a domain event.
func (e *EsdbSerializer) Serialize(
	data domain.IDomainEvent,
	meta metadata.Metadata,
) (*kdb.EventData, error) {
	serializedData, err := e.eventSerializer.Serialize(data)
	if err != nil {
		return nil, err
	}

	serializedMeta, err := e.metadataSerializer.Serialize(meta)
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4()

	googleID, err := googleuuid.Parse(id.String())
	if err != nil {
		return nil, err
	}

	return &kdb.EventData{
		EventID:     googleID,
		EventType:   typeMapper.GetTypeName(data),
		Data:        serializedData.Data,
		ContentType: kdb.ContentTypeJson,
		Metadata:    serializedMeta,
	}, nil
}

// SerializeObject serializes a object.
func (e *EsdbSerializer) SerializeObject(
	data interface{},
	meta metadata.Metadata,
) (*kdb.EventData, error) {
	serializedData, err := e.eventSerializer.SerializeObject(data)
	if err != nil {
		return nil, err
	}

	serializedMeta, err := e.metadataSerializer.Serialize(meta)
	if err != nil {
		return nil, err
	}

	id := uuid.NewV4()

	googleID, err := googleuuid.Parse(id.String())
	if err != nil {
		return nil, err
	}

	return &kdb.EventData{
		EventID:     googleID,
		EventType:   typeMapper.GetTypeName(data),
		Data:        serializedData.Data,
		ContentType: kdb.ContentTypeJson,
		Metadata:    serializedMeta,
	}, nil
}

// Deserialize deserializes a resolved event.
func (e *EsdbSerializer) Deserialize(
	resolveEvent *kdb.ResolvedEvent,
) (domain.IDomainEvent, metadata.Metadata, error) {
	eventType := resolveEvent.Event.EventType
	data := resolveEvent.Event.Data
	userMeta := resolveEvent.Event.UserMetadata

	payload, err := e.eventSerializer.Deserialize(
		data,
		eventType,
		resolveEvent.Event.ContentType,
	)
	if err != nil {
		return nil, nil, err
	}

	meta, err := e.metadataSerializer.Deserialize(userMeta)
	if err != nil {
		return nil, nil, err
	}

	return payload, meta, nil
}

// DeserializeObject deserializes a resolved event.
func (e *EsdbSerializer) DeserializeObject(
	resolveEvent *kdb.ResolvedEvent,
) (interface{}, metadata.Metadata, error) {
	eventType := resolveEvent.Event.EventType
	data := resolveEvent.Event.Data
	userMeta := resolveEvent.Event.UserMetadata

	payload, err := e.eventSerializer.Deserialize(
		data,
		eventType,
		resolveEvent.Event.ContentType,
	)
	if err != nil {
		return nil, nil, err
	}

	meta, err := e.metadataSerializer.Deserialize(userMeta)
	if err != nil {
		return nil, nil, err
	}

	return payload, meta, nil
}

// DomainEventToStreamEvent converts a domain event to a stream event.
func (e *EsdbSerializer) DomainEventToStreamEvent(
	domainEvent domain.IDomainEvent,
	meta metadata.Metadata,
	position int64,
) *models.StreamEvent {
	return &models.StreamEvent{
		EventID:  uuid.NewV4(),
		Event:    domainEvent,
		Metadata: meta,
		Version:  position,
		Position: position,
	}
}
