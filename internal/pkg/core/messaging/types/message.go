// Package types provides message.
package types

import (
	"time"

	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// IMessage is a type that represents a message.
type IMessage interface {
	GeMessageId() string
	GetCreated() time.Time
	// GetMessageTypeName get short type name of the message - we use message short type name instead of full type name because this message in other receiver packages could have different package name
	GetMessageTypeName() string
	GetMessageFullTypeName() string
}

// Message is a struct that represents a message.
type Message struct {
	MessageId string    `json:"messageId,omitempty"`
	Created   time.Time `json:"created"`
	EventType string    `json:"eventType"`
}

// NewMessage is a function that creates a new message.
func NewMessage(messageId string) *Message {
	return &Message{MessageId: messageId, Created: time.Now()}
}

// NewMessageWithTypeName is a function that creates a new message with a type name.
func NewMessageWithTypeName(messageId string, eventTypeName string) *Message {
	return &Message{MessageId: messageId, Created: time.Now(), EventType: eventTypeName}
}

// GeMessageId is a function that returns the message id.
func (m *Message) GeMessageId() string {
	return m.MessageId
}

// GetCreated is a function that returns the created time.
func (m *Message) GetCreated() time.Time {
	return m.Created
}

// GetMessageTypeName is a function that returns the message type name.
func (m *Message) GetMessageTypeName() string {
	return typeMapper.GetTypeName(m)
}

// GetMessageFullTypeName is a function that returns the message full type name.
func (m *Message) GetMessageFullTypeName() string {
	return typeMapper.GetFullTypeName(m)
}
