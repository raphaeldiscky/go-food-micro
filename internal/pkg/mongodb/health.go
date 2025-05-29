package mongodb

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

type mongoHealthChecker struct {
	client *mongo.Client
}

func NewMongoHealthChecker(client *mongo.Client) contracts.Health {
	return &mongoHealthChecker{client}
}

func (healthChecker *mongoHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx, nil)
}

func (healthChecker *mongoHealthChecker) GetHealthName() string {
	return "mongodb"
}
