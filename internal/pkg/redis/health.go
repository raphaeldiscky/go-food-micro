package redis

import (
	"context"

	redis "github.com/redis/go-redis/v9"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

type RedisHealthChecker struct {
	client *redis.Client
}

func NewRedisHealthChecker(client *redis.Client) contracts.Health {
	return &RedisHealthChecker{client}
}

func (healthChecker *RedisHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.Ping(ctx).Err()
}

func (healthChecker *RedisHealthChecker) GetHealthName() string {
	return "redis"
}
