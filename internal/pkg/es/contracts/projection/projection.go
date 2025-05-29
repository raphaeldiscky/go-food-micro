// Package projection provides the projection.
package projection

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

// IProjection is a interface that represents the projection.
type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
