// Package projections contains the elastic order projection.
package projections

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/models"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/orders/contracts/repositories"
)

// elasticOrderProjection is the projection for the order.
type elasticOrderProjection struct {
	elasticOrderReadRepository repositories.OrderElasticRepository
}

// NewElasticOrderProjection creates a new elastic order projection.
func NewElasticOrderProjection(
	elasticOrderReadRepository repositories.OrderElasticRepository,
) projection.IProjection {
	return &elasticOrderProjection{elasticOrderReadRepository: elasticOrderReadRepository}
}

// ProcessEvent processes the event and projects it to the elastic read model.
func (e elasticOrderProjection) ProcessEvent(
	ctx context.Context,
	streamEvent *models.StreamEvent,
) error {
	// TODO: Handling and projecting event to elastic read model
	return nil
}
