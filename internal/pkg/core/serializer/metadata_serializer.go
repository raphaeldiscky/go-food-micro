package serializer

import "github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"

type MetadataSerializer interface {
	Serialize(meta metadata.Metadata) ([]byte, error)
	Deserialize(bytes []byte) (metadata.Metadata, error)
}
