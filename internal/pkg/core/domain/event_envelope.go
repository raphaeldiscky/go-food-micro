package domain

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
