// Package serializer provides a event serialization result.
package serializer

// EventSerializationResult is a struct that represents a event serialization result.
type EventSerializationResult struct {
	Data        []byte
	ContentType string
}
