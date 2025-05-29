// Package serializer provides a metadata serializer.
package serializer

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"

// MetadataSerializer is an interface that represents a metadata serializer.
type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
