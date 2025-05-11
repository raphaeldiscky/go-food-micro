package projection

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"
)

type IProjection interface {
	ProcessEvent(ctx context.Context, streamEvent *models.StreamEvent) error
}
