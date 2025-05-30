// Package tracing provides a message carrier.
package tracing

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"

// MessageCarrier is a struct that represents a message carrier.
type MessageCarrier struct {
	meta *metadata.Metadata
}

// NewMessageCarrier is a function that creates a new message carrier.
func NewMessageCarrier(meta *metadata.Metadata) MessageCarrier {
	return MessageCarrier{meta: meta}
}

// Get is a function that gets a value from the message carrier.
func (a MessageCarrier) Get(key string) string {
	return a.meta.GetString(key)
}

// Set is a function that sets a value in the message carrier.
func (a MessageCarrier) Set(key string, value string) {
	a.meta.Set(key, value)
}

// Keys is a function that returns the keys of the message carrier.
func (a MessageCarrier) Keys() []string {
	return a.meta.Keys()
}
