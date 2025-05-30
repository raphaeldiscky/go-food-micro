// Package types provides message envelope.
package types

// MessageEnvelope is a struct that represents a message envelope.
type MessageEnvelope struct {
	Message IMessage
	Headers map[string]interface{}
}

// NewMessageEnvelope is a function that creates a new message envelope.
func NewMessageEnvelope(
	message IMessage,
	headers map[string]interface{},
) *MessageEnvelope {
	if headers == nil {
		headers = make(map[string]interface{})
	}

	return &MessageEnvelope{
		Message: message,
		Headers: headers,
	}
}

// MessageEnvelopeT is a struct that represents a message envelope.
type MessageEnvelopeT[T IMessage] struct {
	*MessageEnvelope
	Message T
}

// NewMessageEnvelopeT is a function that creates a new message envelope.
func NewMessageEnvelopeT[T IMessage](
	message T,
	headers map[string]interface{},
) *MessageEnvelopeT[T] {
	messageEnvelope := NewMessageEnvelope(message, headers)

	return &MessageEnvelopeT[T]{
		MessageEnvelope: messageEnvelope,
		Message:         message,
	}
}
