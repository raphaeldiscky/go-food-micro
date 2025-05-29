// Package json provides a json event serializer.
package json

import (
	"reflect"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/domain"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// DefaultEventJSONSerializer is a struct that represents a default event json serializer.
type DefaultEventJSONSerializer struct {
	serializer serializer.Serializer
}

// NewDefaultEventJsonSerializer is a function that creates a new default event json serializer.
func NewDefaultEventJsonSerializer(serializer serializer.Serializer) serializer.EventSerializer {
	return &DefaultEventJSONSerializer{serializer: serializer}
}

// Serialize is a function that serializes an event.
func (s *DefaultEventJSONSerializer) Serialize(
	event domain.IDomainEvent,
) (*serializer.EventSerializationResult, error) {
	return s.SerializeObject(event)
}

// SerializeObject is a function that serializes an object.
func (s *DefaultEventJSONSerializer) SerializeObject(
	event interface{},
) (*serializer.EventSerializationResult, error) {
	if event == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: s.ContentType()}, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(event)

	data, err := s.serializer.Marshal(event)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: s.ContentType()}

	return result, nil
}

// Deserialize is a function that deserializes an event.
func (s *DefaultEventJSONSerializer) Deserialize(
	data []byte,
	eventType string,
	contentType string,
) (domain.IDomainEvent, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.EmptyInstanceByTypeNameAndImplementedInterface[domain.IDomainEvent](
		eventType,
	)

	if targetEventPointer == nil {
		return nil, errors.Errorf(
			"event type `%s` is not impelemted IDomainEvent or can't be instansiated",
			eventType,
		)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer.(domain.IDomainEvent), nil
}

// DeserializeObject is a function that deserializes an object.
func (s *DefaultEventJSONSerializer) DeserializeObject(
	data []byte,
	eventType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetEventPointer := typeMapper.InstanceByTypeName(eventType)

	if targetEventPointer == nil {
		return nil, errors.Errorf("event type `%s` can't be instansiated", eventType)
	}

	if contentType != s.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := s.serializer.Unmarshal(data, targetEventPointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", eventType)
	}

	return targetEventPointer, nil
}

// DeserializeType is a function that deserializes a type.
func (s *DefaultEventJSONSerializer) DeserializeType(
	data []byte,
	eventType reflect.Type,
	contentType string,
) (domain.IDomainEvent, error) {
	if data == nil {
		return nil, nil
	}

	// we use event short type name instead of full type name because this event in other receiver packages could have different package name
	eventTypeName := typeMapper.GetTypeName(eventType)

	return s.Deserialize(data, eventTypeName, contentType)
}

// ContentType is a function that returns the content type.
func (s *DefaultEventJSONSerializer) ContentType() string {
	return "application/json"
}

// Serializer is a function that returns the serializer.
func (s *DefaultEventJSONSerializer) Serializer() serializer.Serializer {
	return s.serializer
}
