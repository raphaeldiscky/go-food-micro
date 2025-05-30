// Package redis provides a set of functions for the redis package.
package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// RedisHealthChecker is a struct that contains the redis client.
type RedisHealthChecker struct {
	client *redis.Client
}

// NewRedisHealthChecker creates a new redis health checker.
func NewRedisHealthChecker(client *redis.Client) contracts.Health {
	return &RedisHealthChecker{client}
}

// CheckHealth checks the health of the redis client.
func (healthChecker *RedisHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx).Err()
}

// GetHealthName returns the name of the redis health checker.
func (healthChecker *RedisHealthChecker) GetHealthName() string {
	return "redis"
}
