// Package domain provides a module for the domain.
package domain

import (
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/events"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models/streamversion"
)

// IDomainEvent is a domain event.
type IDomainEvent interface {
	events.IEvent
	GetAggregateId() uuid.UUID
	GetAggregateSequenceNumber() int64
	WithAggregate(aggregateId uuid.UUID, aggregateSequenceNumber int64) *DomainEvent
}

// DomainEvent is a domain event.
type DomainEvent struct {
	*events.Event
	AggregateId             uuid.UUID `json:"aggregate_id"`
	AggregateSequenceNumber int64     `json:"aggregate_sequence_number"`
}

// NewDomainEvent creates a new domain event.
func NewDomainEvent(eventType string) *DomainEvent {
	domainEvent := &DomainEvent{
		Event:                   events.NewEvent(eventType),
		AggregateSequenceNumber: streamversion.NoStream.Value(),
	}
	domainEvent.Event = events.NewEvent(eventType)

	return domainEvent
}

// GetAggregateId gets the aggregate id.
func (d *DomainEvent) GetAggregateId() uuid.UUID {
	return d.AggregateId
}

// GetAggregateSequenceNumber gets the aggregate sequence number.
func (d *DomainEvent) GetAggregateSequenceNumber() int64 {
	return d.AggregateSequenceNumber
}

// WithAggregate sets the aggregate id and sequence number.
func (d *DomainEvent) WithAggregate(
	aggregateId uuid.UUID,
	aggregateSequenceNumber int64,
) *DomainEvent {
	d.AggregateId = aggregateId
	d.AggregateSequenceNumber = aggregateSequenceNumber

	return d
}
