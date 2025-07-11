// Package models provides a event sourced aggregate.
package models

// https://www.eventstore.com/blog/what-is-event-sourcing
// https://www.eventstore.com/blog/event-sourcing-and-cqrs

import (
	"encoding/json"
	"fmt"

	"emperror.dev/errors"

	linq "github.com/ahmetb/go-linq/v3"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	errors2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/es/errors"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
)

// WhenFunc is a function that updates the aggregate state with new events that are added to the event store and also for events that are already in the event store without increasing the version.
type WhenFunc func(event domain.IDomainEvent) error

// When is a interface that updates the aggregate state with new events that are added to the event store and also for events that are already in the event store without increasing the version.
type When interface {
	// When Update the aggregate state with new events that are added to the event store and also for events that are already in the event store without increasing the version.
	When(event domain.IDomainEvent) error
}

// fold is a interface that restores the aggregate state with events that are loaded form the event store and increase the current version and last commit version.
type fold interface {
	// Restore the aggregate state with events that are loaded form the event store and increase the current version and last commit version.
	fold(event domain.IDomainEvent, metadata metadata.Metadata) error
}

// Apply is a interface that applies a new event to the aggregate state, adds the event to the list of pending changes,
// and increases the `CurrentVersion` property and `LastCommittedVersion` will be unchanged.
type Apply interface {
	// Apply a new event to the aggregate state, adds the event to the list of pending changes,
	// and increases the `CurrentVersion` property and `LastCommittedVersion` will be unchanged.
	Apply(event domain.IDomainEvent, isNew bool) error
}

// AggregateStateProjection is a interface that applies a new event to the aggregate state, adds the event to the list of pending changes,
// and increases the `CurrentVersion` property and `LastCommittedVersion` will be unchanged.
type AggregateStateProjection interface {
	Apply
	fold
}

// IHaveEventSourcedAggregate is a interface that should implement by actual aggregate root class in our domain_events.
type IHaveEventSourcedAggregate interface {
	When
	NewEmptyAggregate()
	IEventSourcedAggregateRoot
}

// IEventSourcedAggregateRoot is a interface that contains all methods of AggregateBase.
type IEventSourcedAggregateRoot interface {
	domain.IEntity

	// OriginalVersion Gets the original version is the aggregate version we got from the store. This is used to ensure optimistic concurrency,
	// to check if there were no changes made to the aggregate state between load and save for the current operation.
	OriginalVersion() int64

	SetOriginalVersion(version int64)

	// CurrentVersion Gets the current version is set to original version when the aggregate is loaded from the store.
	// It should increase for each state transition performed within the scope of the current operation.
	CurrentVersion() int64

	// AddDomainEvents adds a new domain_events event to the aggregate's uncommitted events.
	AddDomainEvents(event domain.IDomainEvent) error

	// MarkUncommittedEventAsCommitted Mark all changes (events) as committed, clears uncommitted changes and updates the current version of the aggregate.
	MarkUncommittedEventAsCommitted()

	// HasUncommittedEvents Does the aggregate have change that have not been committed to storage
	HasUncommittedEvents() bool

	// UncommittedEvents Gets a list of uncommitted events for this aggregate.
	UncommittedEvents() []domain.IDomainEvent

	// LoadFromHistory Loads the current state of the aggregate from a list of events.
	LoadFromHistory(events []domain.IDomainEvent, metadata metadata.Metadata) error

	AggregateStateProjection
}

// EventSourcedAggregateRoot base aggregate contains all main necessary fields.
type EventSourcedAggregateRoot struct {
	*domain.Entity
	originalVersion   int64
	currentVersion    int64
	uncommittedEvents []domain.IDomainEvent
	when              WhenFunc
}

// EventSourcedAggregateRootDataModel is a data model for the event sourced aggregate root.
type EventSourcedAggregateRootDataModel struct {
	*domain.EntityDataModel
	OriginalVersion int64 `json:"originalVersion" bson:"originalVersion"`
}

// NewEventSourcedAggregateRootWithID creates a new event sourced aggregate root with an id.
func NewEventSourcedAggregateRootWithID(
	id uuid.UUID,
	aggregateType string,
	when WhenFunc,
) *EventSourcedAggregateRoot {
	if when == nil {
		return nil
	}

	aggregate := &EventSourcedAggregateRoot{
		originalVersion: streamversion.NoStream.Value(),
		currentVersion:  streamversion.NoStream.Value(),
		when:            when,
	}

	aggregate.Entity = domain.NewEntityWithID(id, aggregateType)

	return aggregate
}

// NewEventSourcedAggregateRoot creates a new event sourced aggregate root.
func NewEventSourcedAggregateRoot(aggregateType string, when WhenFunc) *EventSourcedAggregateRoot {
	if when == nil {
		return nil
	}

	aggregate := &EventSourcedAggregateRoot{
		originalVersion: streamversion.NoStream.Value(),
		currentVersion:  streamversion.NoStream.Value(),
		when:            when,
	}

	aggregate.Entity = domain.NewEntity(aggregateType)

	return aggregate
}

// OriginalVersion gets the original version is the aggregate version we got from the store. This is used to ensure optimistic concurrency,
// to check if there were no changes made to the aggregate state between load and save for the current operation.
func (a *EventSourcedAggregateRoot) OriginalVersion() int64 {
	return a.originalVersion
}

// SetOriginalVersion sets the original version.
func (a *EventSourcedAggregateRoot) SetOriginalVersion(version int64) {
	a.originalVersion = version
}

// CurrentVersion gets the current version is set to original version when the aggregate is loaded from the store.
// It should increase for each state transition performed within the scope of the current operation.
func (a *EventSourcedAggregateRoot) CurrentVersion() int64 {
	return a.currentVersion
}

// AddDomainEvents adds a new domain_events event to the aggregate's uncommitted events.
func (a *EventSourcedAggregateRoot) AddDomainEvents(event domain.IDomainEvent) error {
	exists := linq.From(a.uncommittedEvents).AnyWithT(func(e domain.IDomainEvent) bool {
		return e.GetEventID() == event.GetEventID()
	})

	if exists {
		return errors2.EventAlreadyExistsError
	}
	event.WithAggregate(a.ID(), a.CurrentVersion()+1)
	a.uncommittedEvents = append(a.uncommittedEvents, event)

	return nil
}

// MarkUncommittedEventAsCommitted Mark all changes (events) as committed, clears uncommitted changes and updates the current version of the aggregate.
func (a *EventSourcedAggregateRoot) MarkUncommittedEventAsCommitted() {
	a.uncommittedEvents = nil
}

// HasUncommittedEvents Does the aggregate have change that have not been committed to storage.
func (a *EventSourcedAggregateRoot) HasUncommittedEvents() bool {
	return len(a.uncommittedEvents) > 0
}

// UncommittedEvents Gets a list of uncommitted events for this aggregate.
func (a *EventSourcedAggregateRoot) UncommittedEvents() []domain.IDomainEvent {
	return a.uncommittedEvents
}

// LoadFromHistory Loads the current state of the aggregate from a list of events.
func (a *EventSourcedAggregateRoot) LoadFromHistory(
	events []domain.IDomainEvent,
	metadata metadata.Metadata,
) error {
	for _, event := range events {
		err := a.fold(event, metadata)
		if err != nil {
			return errors.WrapIf(
				err,
				"[EventSourcedAggregateRoot_LoadFromHistory:fold] error in loading event from history",
			)
		}
	}

	return nil
}

// Apply applies a new event to the aggregate state, adds the event to the list of pending changes,
// and increases the `CurrentVersion` property and `LastCommittedVersion` will be unchanged.
func (a *EventSourcedAggregateRoot) Apply(event domain.IDomainEvent, isNew bool) error {
	if isNew {
		err := a.AddDomainEvents(event)
		if err != nil {
			return errors.WrapIf(
				err,
				"[EventSourcedAggregateRoot_Apply:AddDomainEvents] error in adding domain_events event to the domain_events events list",
			)
		}
	}
	err := a.when(event)
	if err != nil {
		return errors.WrapIf(err, "[EventSourcedAggregateRoot_Apply:when] error in the whenFunc")
	}
	a.currentVersion++

	return nil
}

// fold restores the aggregate state with events that are loaded form the event store and increase the current version and last commit version.
func (a *EventSourcedAggregateRoot) fold(
	event domain.IDomainEvent,
	_ metadata.Metadata,
) error {
	err := a.when(event)
	if err != nil {
		return errors.WrapIf(
			err,
			"[EventSourcedAggregateRoot_fold:when] error in the applying whenFunc",
		)
	}
	a.originalVersion++
	a.currentVersion++

	return nil
}

// String returns a string representation of the event sourced aggregate root.
func (a *EventSourcedAggregateRoot) String() string {
	data := &EventSourcedAggregateRootDataModel{
		EntityDataModel: a.ToDataModel(),
		OriginalVersion: a.originalVersion,
	}
	j, err := json.Marshal(data)
	if err != nil {
		return fmt.Sprintf("Aggregate json: %v", err)
	}

	return fmt.Sprintf("Aggregate json: %s", string(j))
}
