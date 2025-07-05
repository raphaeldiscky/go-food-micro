// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"context"
	"fmt"
	"io"
	"time"

	"emperror.dev/errors"

	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/events"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// esdbSubscriptionCheckpointRepository is a struct that represents a event store db subscription checkpoint repository.
type esdbSubscriptionCheckpointRepository struct {
	client        *kdb.Client
	log           logger.Logger
	esdbSerilizer *EsdbSerializer
}

// CheckpointStored is a struct that represents a checkpoint stored.
type CheckpointStored struct {
	Position       uint64
	SubscriptionId string
	CheckpointAt   time.Time
	*events.Event
}

// NewEsdbSubscriptionCheckpointRepository creates a new event store db subscription checkpoint repository.
func NewEsdbSubscriptionCheckpointRepository(
	client *kdb.Client,
	logger logger.Logger,
	esdbSerializer *EsdbSerializer,
) contracts.SubscriptionCheckpointRepository {
	return &esdbSubscriptionCheckpointRepository{
		client:        client,
		log:           logger,
		esdbSerilizer: esdbSerializer,
	}
}

// Load loads a subscription checkpoint.
func (e *esdbSubscriptionCheckpointRepository) Load(
	subscriptionId string,
	ctx context.Context,
) (uint64, error) {
	streamName := getCheckpointStreamName(subscriptionId)

	stream, err := e.client.ReadStream(
		ctx,
		streamName,
		kdb.ReadStreamOptions{
			Direction: kdb.Backwards,
			From:      kdb.End{},
		}, 1)

	var kdbErr *kdb.Error
	if errors.As(err, &kdbErr) && kdbErr.Code() == kdb.ErrorCodeResourceNotFound {
		return 0, nil
	} else if err != nil {
		return 0, errors.WrapIf(err, "db.ReadStream")
	}

	event, err := stream.Recv()
	if errors.As(err, &kdbErr) && kdbErr.Code() == kdb.ErrorCodeResourceNotFound {
		return 0, nil
	}
	if errors.Is(err, io.EOF) {
		return 0, nil
	}
	if err != nil {
		return 0, errors.WrapIf(err, "stream.Recv")
	}

	deserialized, _, err := e.esdbSerilizer.DeserializeObject(event)
	if err != nil {
		return 0, err
	}

	v, ok := deserialized.(*CheckpointStored)
	if !ok {
		return 0, nil
	}

	stream.Close()

	return v.Position, nil
}

// Store stores a subscription checkpoint.
func (e *esdbSubscriptionCheckpointRepository) Store(
	subscriptionID string,
	position uint64,
	ctx context.Context,
) error {
	checkpoint := &CheckpointStored{
		SubscriptionId: subscriptionID,
		Position:       position,
		CheckpointAt:   time.Now(),
		Event:          events.NewEvent(typeMapper.GetTypeName(&CheckpointStored{})),
	}
	streamName := getCheckpointStreamName(subscriptionID)
	eventData, err := e.esdbSerilizer.SerializeObject(checkpoint, nil)
	if err != nil {
		return errors.WrapIf(err, "esdbSerilizer.Serialize")
	}

	// First, try to append to stream assuming it exists
	_, err = e.client.AppendToStream(
		ctx,
		streamName,
		kdb.AppendToStreamOptions{StreamState: kdb.StreamExists{}},
		*eventData,
	)

	var wrongVersionErr *kdb.Error
	if errors.As(err, &wrongVersionErr) &&
		wrongVersionErr.Code() == kdb.ErrorCodeWrongExpectedVersion {
		// Stream doesn't exist, so we need to create it
		// First, set the stream metadata to have at most 1 event
		streamMeta := kdb.StreamMetadata{}
		streamMeta.SetMaxCount(1)

		_, err := e.client.SetStreamMetadata(
			ctx,
			streamName,
			kdb.AppendToStreamOptions{StreamState: kdb.NoStream{}},
			streamMeta)
		if err != nil {
			return errors.WrapIf(err, "client.SetStreamMetadata")
		}

		// Now append the first event to the new stream
		_, err = e.client.AppendToStream(
			ctx,
			streamName,
			kdb.AppendToStreamOptions{StreamState: kdb.NoStream{}},
			*eventData,
		)
		if err != nil {
			return errors.WrapIf(err, "client.AppendToStream after metadata set")
		}
	} else if err != nil {
		return errors.WrapIf(err, "client.AppendToStream")
	}

	return nil
}

// getCheckpointStreamName gets a checkpoint stream name.
func getCheckpointStreamName(subscriptionId string) string {
	return fmt.Sprintf("$checkpoint_stream_%s", subscriptionId)
}
