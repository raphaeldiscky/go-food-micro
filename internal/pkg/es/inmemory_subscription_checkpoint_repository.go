//go:build unit
// +build unit

// Package es provides a in memory subscription checkpoint repository.
package es

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts"
)

// inMemorySubscriptionCheckpointRepository is a in memory subscription checkpoint repository.
type inMemorySubscriptionCheckpointRepository struct {
	checkpoints map[string]uint64
}

// NewInMemorySubscriptionCheckpointRepository creates a new in memory subscription checkpoint repository.
func NewInMemorySubscriptionCheckpointRepository() contracts.SubscriptionCheckpointRepository {
	return &inMemorySubscriptionCheckpointRepository{checkpoints: make(map[string]uint64)}
}

// Load loads a subscription checkpoint.
func (i inMemorySubscriptionCheckpointRepository) Load(
	subscriptionId string,
	ctx context.Context,
) (uint64, error) {
	checkpoint := i.checkpoints[subscriptionId]
	if checkpoint == 0 {
		return 0, nil
	}

	return checkpoint, nil
}

// Store stores a subscription checkpoint.
func (i inMemorySubscriptionCheckpointRepository) Store(
	subscriptionId string,
	position uint64,
	ctx context.Context,
) error {
	i.checkpoints[subscriptionId] = position

	return nil
}
