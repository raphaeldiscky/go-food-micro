package projection

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

type IProjectionPublisher interface {
	Publish(ctx context.Context, streamEvent *models.StreamEvent) error
}
