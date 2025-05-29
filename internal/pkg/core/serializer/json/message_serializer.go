// Package json provides a json message serializer.
package json

import (
	"reflect"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/messaging/types"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// DefaultMessageJsonSerializer is a struct that represents a default message json serializer.
type DefaultMessageJsonSerializer struct {
	serializer serializer.Serializer
}

// NewDefaultMessageJsonSerializer is a function that creates a new default message json serializer.
func NewDefaultMessageJsonSerializer(s serializer.Serializer) serializer.MessageSerializer {
	return &DefaultMessageJsonSerializer{serializer: s}
}

// Serialize is a function that serializes a message.
func (m *DefaultMessageJsonSerializer) Serialize(
	message types.IMessage,
) (*serializer.EventSerializationResult, error) {
	return m.SerializeObject(message)
}

// SerializeObject is a function that serializes an object.
func (m *DefaultMessageJsonSerializer) SerializeObject(
	message interface{},
) (*serializer.EventSerializationResult, error) {
	if message == nil {
		return &serializer.EventSerializationResult{Data: nil, ContentType: m.ContentType()}, nil
	}

	// we use message short type name instead of full type name because this message in other receiver packages could have different package name
	eventType := typeMapper.GetTypeName(message)

	data, err := m.serializer.Marshal(message)
	if err != nil {
		return nil, errors.WrapIff(err, "error in Marshaling: `%s`", eventType)
	}

	result := &serializer.EventSerializationResult{Data: data, ContentType: m.ContentType()}

	return result, nil
}

// SerializeEnvelop is a function that serializes a message envelop.
func (m *DefaultMessageJsonSerializer) SerializeEnvelop(
	messageEnvelop types.MessageEnvelope,
) (*serializer.EventSerializationResult, error) {
	// TODO implement me
	panic("implement me")
}

// Deserialize is a function that deserializes a message.
func (m *DefaultMessageJsonSerializer) Deserialize(
	data []byte,
	messageType string,
	contentType string,
) (types.IMessage, error) {
	if data == nil {
		return nil, nil
	}

	targetMessagePointer := typeMapper.EmptyInstanceByTypeNameAndImplementedInterface[types.IMessage](
		messageType,
	)

	if targetMessagePointer == nil {
		return nil, errors.Errorf(
			"message type `%s` is not impelemted IMessage or can't be instansiated",
			messageType,
		)
	}

	if contentType != m.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := m.serializer.Unmarshal(data, targetMessagePointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", messageType)
	}

	return targetMessagePointer.(types.IMessage), nil
}

// DeserializeObject is a function that deserializes an object.
func (m *DefaultMessageJsonSerializer) DeserializeObject(
	data []byte,
	messageType string,
	contentType string,
) (interface{}, error) {
	if data == nil {
		return nil, nil
	}

	targetMessagePointer := typeMapper.InstanceByTypeName(messageType)

	if targetMessagePointer == nil {
		return nil, errors.Errorf("message type `%s` can't be instansiated", messageType)
	}

	if contentType != m.ContentType() {
		return nil, errors.Errorf("contentType: %s is not supported", contentType)
	}

	if err := m.serializer.Unmarshal(data, targetMessagePointer); err != nil {
		return nil, errors.WrapIff(err, "error in Unmarshaling: `%s`", messageType)
	}

	return targetMessagePointer, nil
}

// DeserializeType is a function that deserializes a type.
func (m *DefaultMessageJsonSerializer) DeserializeType(
	data []byte,
	messageType reflect.Type,
	contentType string,
) (types.IMessage, error) {
	if data == nil {
		return nil, nil
	}

	// we use message short type name instead of full type name because this message in other receiver packages could have different package name
	messageTypeName := typeMapper.GetTypeName(messageType)

	return m.Deserialize(data, messageTypeName, contentType)
}

// ContentType is a function that returns the content type.
func (m *DefaultMessageJsonSerializer) ContentType() string {
	return "application/json"
}

// Serializer is a function that returns the serializer.
func (m *DefaultMessageJsonSerializer) Serializer() serializer.Serializer {
	return m.serializer
}
