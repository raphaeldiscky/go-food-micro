// Package types provides message consume context.
package types

import (
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

// MessageConsumeContext is a type that represents a message consume context.
type MessageConsumeContext interface {
	MessageId() string
	CorrelationId() string
	MessageType() string
	Created() time.Time
	ContentType() string
	DeliveryTag() uint64
	Metadata() metadata.Metadata
	Message() IMessage
}

// messageConsumeContext is a struct that represents a message consume context.
type messageConsumeContext struct {
	metadata      metadata.Metadata
	contentType   string
	messageType   string
	messageId     string
	created       time.Time
	tag           uint64
	correlationId string
	message       IMessage
}

// NewMessageConsumeContext is a function that creates a new message consume context.
func NewMessageConsumeContext(
	message IMessage,
	meta metadata.Metadata,
	contentType string,
	messageType string,
	created time.Time,
	deliveryTag uint64,
	messageId string,
	correlationId string,
) MessageConsumeContext {
	return &messageConsumeContext{
		message:       message,
		metadata:      meta,
		contentType:   contentType,
		messageId:     messageId,
		tag:           deliveryTag,
		created:       created,
		messageType:   messageType,
		correlationId: correlationId,
	}
}

// Message is a function that returns the message.
func (m *messageConsumeContext) Message() IMessage {
	return m.message
}

// MessageId is a function that returns the message id.
func (m *messageConsumeContext) MessageId() string {
	return m.messageId
}

// CorrelationId is a function that returns the correlation id.
func (m *messageConsumeContext) CorrelationId() string {
	return m.correlationId
}

// MessageType is a function that returns the message type.
func (m *messageConsumeContext) MessageType() string {
	return m.messageType
}

// ContentType is a function that returns the content type.
func (m *messageConsumeContext) ContentType() string {
	return m.contentType
}

// Metadata is a function that returns the metadata.
func (m *messageConsumeContext) Metadata() metadata.Metadata {
	return m.metadata
}

// Created is a function that returns the created time.
func (m *messageConsumeContext) Created() time.Time {
	return m.created
}

// DeliveryTag is a function that returns the delivery tag.
func (m *messageConsumeContext) DeliveryTag() uint64 {
	return m.tag
}
