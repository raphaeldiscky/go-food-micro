// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/EventStore/EventStore-Client-Go/esdb"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// esdbSubscriptionAllWorker is a struct that represents a event store db subscription all worker.
type esdbSubscriptionAllWorker struct {
	db                               *esdb.Client
	cfg                              *config.EventStoreDbOptions
	log                              logger.Logger
	subscriptionOption               *EventStoreDBSubscriptionToAllOptions
	esdbSerializer                   *EsdbSerializer
	subscriptionCheckpointRepository contracts.SubscriptionCheckpointRepository
	subscriptionId                   string
	projectionPublisher              projection.IProjectionPublisher
}

// EsdbSubscriptionAllWorker is a interface that represents a event store db subscription all worker.
type EsdbSubscriptionAllWorker interface {
	SubscribeAll(
		ctx context.Context,
		subscriptionOption *EventStoreDBSubscriptionToAllOptions,
	) error
}

// EventStoreDBSubscriptionToAllOptions is a struct that represents a event store db subscription to all options.
type EventStoreDBSubscriptionToAllOptions struct {
	SubscriptionId              string
	FilterOptions               *esdb.SubscriptionFilter
	Credentials                 *esdb.Credentials
	ResolveLinkTos              bool
	IgnoreDeserializationErrors bool
	Prefix                      string
}

// NewEsdbSubscriptionAllWorker creates a new event store db subscription all worker.
func NewEsdbSubscriptionAllWorker(
	log logger.Logger,
	db *esdb.Client,
	cfg *config.EventStoreDbOptions,
	esdbSerializer *EsdbSerializer,
	subscriptionRepository contracts.SubscriptionCheckpointRepository,
	projectionBuilderFunc ProjectionBuilderFuc,
) EsdbSubscriptionAllWorker {
	builder := NewProjectionsBuilder()
	if projectionBuilderFunc != nil {
		projectionBuilderFunc(builder)
	}
	projectionConfigurations := builder.Build()
	projectionPublisher := es.NewProjectionPublisher(projectionConfigurations.Projections)

	return &esdbSubscriptionAllWorker{
		db:                               db,
		cfg:                              cfg,
		log:                              log,
		esdbSerializer:                   esdbSerializer,
		subscriptionCheckpointRepository: subscriptionRepository,
		projectionPublisher:              projectionPublisher,
	}
}

// handleSubscriptionEvent processes a single subscription event and updates the position.
func (s *esdbSubscriptionAllWorker) handleSubscriptionEvent(
	ctx context.Context,
	event *esdb.SubscriptionEvent,
	options *esdb.SubscribeToAllOptions,
) error {
	if event.SubscriptionDropped != nil {
		s.log.Errorf(
			"subscription to all '%s' dropped: %s",
			s.subscriptionId,
			event.SubscriptionDropped.Error,
		)

		return event.SubscriptionDropped.Error
	}

	if event.EventAppeared != nil {
		streamID := event.EventAppeared.OriginalEvent().StreamID
		revision := event.EventAppeared.OriginalEvent().EventNumber
		s.log.Info(
			fmt.Sprintf(
				"event appeared in subscription to all '%s'. streamId: %s, revision: %d",
				s.subscriptionId,
				streamID,
				revision,
			),
		)

		options.From = event.EventAppeared.OriginalEvent().Position

		return s.handleEvent(ctx, event.EventAppeared)
	}

	return nil
}

// SubscribeAll subscribes to all.
func (s *esdbSubscriptionAllWorker) SubscribeAll(
	ctx context.Context,
	subscriptionOption *EventStoreDBSubscriptionToAllOptions,
) error {
	if subscriptionOption.SubscriptionId == "" {
		subscriptionOption.SubscriptionId = "defaultLogger"
	}

	if subscriptionOption.FilterOptions == nil {
		subscriptionOption.FilterOptions = esdb.ExcludeSystemEventsFilter()
	}

	s.subscriptionOption = subscriptionOption
	s.subscriptionId = subscriptionOption.SubscriptionId

	s.log.Info(fmt.Sprintf("starting subscription to all '%s'.", subscriptionOption.SubscriptionId))

	checkpoint, err := s.subscriptionCheckpointRepository.Load(
		subscriptionOption.SubscriptionId,
		ctx,
	)
	if err != nil {
		return err
	}

	var from esdb.AllPosition
	if checkpoint == 0 {
		from = esdb.Start{}
	} else {
		from = esdb.Position{
			Commit:  checkpoint,
			Prepare: checkpoint,
		}
	}

	options := esdb.SubscribeToAllOptions{
		ResolveLinkTos:     subscriptionOption.ResolveLinkTos,
		Authenticated:      subscriptionOption.Credentials,
		Filter:             subscriptionOption.FilterOptions,
		From:               from,
		CheckpointInterval: 1,
	}

	for {
		stream, err := s.db.SubscribeToAll(ctx, options)
		if err != nil {
			time.Sleep(1 * time.Second)

			continue
		}

		s.log.Info(
			fmt.Sprintf("subscription to all '%s' started.", subscriptionOption.SubscriptionId),
		)

		for {
			event := stream.Recv()
			if err := s.handleSubscriptionEvent(ctx, event, &options); err != nil {
				if err := stream.Close(); err != nil {
					s.log.Errorf("error closing stream: %v", err)
				}

				return err
			}
		}

		return ctx.Err()
	}
}

// handleEvent handles an event.
func (s *esdbSubscriptionAllWorker) handleEvent(
	ctx context.Context,
	resolvedEvent *esdb.ResolvedEvent,
) error {
	if s.isCheckpointEvent(resolvedEvent) || s.isEventWithEmptyData(resolvedEvent) {
		return nil
	}

	streamEvent, err := s.esdbSerializer.ResolvedEventToStreamEvent(resolvedEvent)
	if err != nil {
		return errors.WrapIf(err, "failed to convert resolved event to stream event")
	}

	// publish to internal event bus - for handling event and project it manually tp corresponding read model
	err = mediatr.Publish(ctx, streamEvent)
	if err != nil {
		return errors.WrapIf(
			err,
			"failed to publish stream event for the mediatr (internal event bus for handling event)",
		)
	}

	// publish to projection publisher
	err = s.projectionPublisher.Publish(ctx, streamEvent)
	if err != nil {
		return errors.WrapIf(err, "failed to publish stream event in the handle event")
	}

	err = s.subscriptionCheckpointRepository.Store(
		s.subscriptionId,
		resolvedEvent.Event.Position.Commit,
		ctx,
	)
	if err != nil {
		return errors.WrapIf(err, "failed to store subscription checkpoint")
	}

	return nil
}

// isEventWithEmptyData checks if an event has empty data.
func (s *esdbSubscriptionAllWorker) isEventWithEmptyData(resolvedEvent *esdb.ResolvedEvent) bool {
	if len(resolvedEvent.Event.Data) != 0 {
		return false
	}

	s.log.Info("event with empty data received")

	return true
}

// isCheckpointEvent checks if an event is a checkpoint event.
func (s *esdbSubscriptionAllWorker) isCheckpointEvent(resolvedEvent *esdb.ResolvedEvent) bool {
	name := typeMapper.GetFullTypeName(CheckpointStored{})
	if resolvedEvent.Event.EventType != name {
		return false
	}

	s.log.Info("checkpoint event received - skipping")

	return true
}

//https://developers.eventstore.com/clients/grpc/subscriptions.html#handling-subscription-drops
// func (s *esdbSubscriptionAllWorker) resubscribe(ctx context.Context) {
//	for true {
//		err := s.SubscribeAll(ctx, s.subscriptionOption)
//		if err != nil {
//			s.log.Error(err)
//			time.Sleep(time.Second * 1000)
//			continue
//		}
//
//		break
//	}
//}
