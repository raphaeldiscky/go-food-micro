// Package projection provides the projection publisher.
package projection

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

// IProjectionPublisher is a interface that represents the projection publisher.
type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
