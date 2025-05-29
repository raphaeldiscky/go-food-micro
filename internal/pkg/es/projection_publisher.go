// Package es provides a projection publisher.
package es

import (
	"context"

	"emperror.dev/errors"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

// projectionPublisher is a projection publisher.
type projectionPublisher struct {
	projections []projection.IProjection
}

// NewProjectionPublisher creates a new projection publisher.
func NewProjectionPublisher(projections []projection.IProjection) projection.IProjectionPublisher {
	return &projectionPublisher{projections: projections}
}

// Publish publishes a stream event.
func (p projectionPublisher) Publish(ctx context.Context, streamEvent *models.StreamEvent) error {
	if streamEvent == nil {
		return nil
	}

	if p.projections == nil {
		return nil
	}

	for _, pj := range p.projections {
		err := pj.ProcessEvent(ctx, streamEvent)
		if err != nil {
			return errors.WrapIf(err, "error in processing projection")
		}
	}

	return nil
}
