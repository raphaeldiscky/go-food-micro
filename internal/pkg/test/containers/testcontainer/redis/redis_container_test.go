//go:build integration
// +build integration

// Package redis provides a redis container.
package redis

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	redis "github.com/redis/go-redis/v9"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	redis2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
)

// TestCustomRedisContainer tests the custom redis container.
func TestCustomRedisContainer(t *testing.T) {
	ctx := context.Background()
	var redisClient redis.UniversalClient

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		redis2.Module,
		fx.Decorate(RedisContainerOptionsDecorator(t, ctx)),
		fx.Populate(&redisClient),
	).RequireStart()

	assert.NotNil(t, redisClient)
}
