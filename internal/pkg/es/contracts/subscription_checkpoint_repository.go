// Package contracts provides a set of contracts for the es package.
package contracts

import "context"

// SubscriptionCheckpointRepository is a repository for subscription checkpoints.
type SubscriptionCheckpointRepository interface {
	Load(subscriptionID string, ctx context.Context) (uint64, error)
	Store(subscriptionID string, position uint64, ctx context.Context) error
}
