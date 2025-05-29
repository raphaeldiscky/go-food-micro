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

func Test_Custom_Redis_Container(t *testing.T) {
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
