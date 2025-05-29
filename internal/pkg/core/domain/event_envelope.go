// Package domain provides a module for the domain.
package domain

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core/metadata"
)

// EventEnvelope is a event envelope.
type EventEnvelope struct {
	EventData interface{}
	Metadata  metadata.Metadata
}
