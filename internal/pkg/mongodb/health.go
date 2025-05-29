// Package mongodb provides a health checker for the mongodb.
package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// mongoHealthChecker is a health checker for the mongodb.
type mongoHealthChecker struct {
	client *mongo.Client
}

// NewMongoHealthChecker creates a new mongodb health checker.
func NewMongoHealthChecker(client *mongo.Client) contracts.Health {
	return &mongoHealthChecker{client}
}

// CheckHealth checks the health of the mongodb.
func (healthChecker *mongoHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx, nil)
}

// GetHealthName gets the health name.
func (healthChecker *mongoHealthChecker) GetHealthName() string {
	return "mongodb"
}
