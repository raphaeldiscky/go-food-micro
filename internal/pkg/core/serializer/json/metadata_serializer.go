// Package json provides a json metadata serializer.
package json

import (
	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/serializer"
)

// DefaultMetadataJsonSerializer is a struct that represents a default metadata json serializer.
type DefaultMetadataJsonSerializer struct {
	serializer serializer.Serializer
}

// NewDefaultMetadataJsonSerializer is a function that creates a new default metadata json serializer.
func NewDefaultMetadataJsonSerializer(
	serializer serializer.Serializer,
) serializer.MetadataSerializer {
	return &DefaultMetadataJsonSerializer{serializer: serializer}
}

// Serialize is a function that serializes a metadata.
func (s *DefaultMetadataJsonSerializer) Serialize(meta metadata.Metadata) ([]byte, error) {
	if meta == nil {
		return nil, nil
	}

	marshal, err := s.serializer.Marshal(meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to marshal metadata")
	}

	return marshal, nil
}

// Deserialize is a function that deserializes a metadata.
func (s *DefaultMetadataJsonSerializer) Deserialize(
	bytes []byte,
) (metadata.Metadata, error) {
	if bytes == nil {
		return nil, nil
	}

	var meta metadata.Metadata

	err := s.serializer.Unmarshal(bytes, &meta)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to unmarshal metadata")
	}

	return meta, nil
}
