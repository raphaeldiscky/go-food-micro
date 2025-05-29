// Package events provides a module for the events.
package events

import (
	"time"

	uuid "github.com/satori/go.uuid"

	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// IEvent is a event.
type IEvent interface {
	GetEventID() uuid.UUID
	GetOccurredOn() time.Time
	// GetEventTypeName get short type name of the event - we use event short type name instead of full type name because this event in other receiver packages could have different package name
	GetEventTypeName() string
	GetEventFullTypeName() string
}

// Event is a event.
type Event struct {
	EventId    uuid.UUID `json:"event_id"`
	EventType  string    `json:"event_type"`
	OccurredOn time.Time `json:"occurred_on"`
}

// NewEvent creates a new event.
func NewEvent(eventType string) *Event {
	return &Event{
		EventId:    uuid.NewV4(),
		OccurredOn: time.Now(),
		EventType:  eventType,
	}
}

// GetEventID gets the event id.
func (e *Event) GetEventID() uuid.UUID {
	return e.EventId
}

// GetEventType gets the event type.
func (e *Event) GetEventType() string {
	return e.EventType
}

// GetOccurredOn gets the occurred on.
func (e *Event) GetOccurredOn() time.Time {
	return e.OccurredOn
}

// GetEventTypeName gets the event type name.
func (e *Event) GetEventTypeName() string {
	return typeMapper.GetTypeName(e)
}

// GetEventFullTypeName gets the event full type name.
func (e *Event) GetEventFullTypeName() string {
	return typeMapper.GetFullTypeName(e)
}

// IsEvent checks if the object is a event.
func IsEvent(obj interface{}) bool {
	if _, ok := obj.(IEvent); ok {
		return true
	}

	return false
}
