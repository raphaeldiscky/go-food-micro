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

	_, err = e.client.AppendToStream(
		ctx,
		streamName,
		kdb.AppendToStreamOptions{StreamState: kdb.StreamExists{}},
		*eventData,
	)

	var wrongVersionErr *kdb.Error
	if errors.As(err, &wrongVersionErr) &&
		wrongVersionErr.Code() == kdb.ErrorCodeWrongExpectedVersion {
		streamMeta := kdb.StreamMetadata{}
		streamMeta.SetMaxCount(1)

		// WrongExpectedVersionException means that stream did not exist
		// Set the checkpoint stream to have at most 1 event
		// using stream metadata $maxCount property
		_, err := e.client.SetStreamMetadata(
			ctx,
			streamName,
			kdb.AppendToStreamOptions{StreamState: kdb.StreamExists{}},
			streamMeta)
		if err != nil {
			return errors.WrapIf(err, "client.SetStreamMetadata")
		}

		// append event again expecting stream to not exist
		_, err = e.client.AppendToStream(
			ctx,
			streamName,
			kdb.AppendToStreamOptions{StreamState: kdb.StreamExists{}},
			*eventData,
		)
		if err != nil {
			return err
		}
	} else {
		return err
	}

	return nil
}

// getCheckpointStreamName gets a checkpoint stream name.
func getCheckpointStreamName(subscriptionId string) string {
	return fmt.Sprintf("$checkpoint_stream_%s", subscriptionId)
}
