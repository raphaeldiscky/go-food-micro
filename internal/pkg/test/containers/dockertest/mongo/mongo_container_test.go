//go:build integration
// +build integration

// Package mongo provides the mongo container test.
package mongo

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
)

// TestMongoContainer tests the mongo container.
func TestMongoContainer(t *testing.T) {
	t.Skip(
		"Skipping Mongo dockertest container test due to persistent port conflicts. See issue with dockertest port allocation.",
	)
	ctx := context.Background()
	var mongoClient *mongo.Client

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		mongodb.Module,
		fx.Decorate(MongoDockerTestContainerOptionsDecorator(t, ctx)),
		fx.Populate(&mongoClient),
	).RequireStart()

	assert.NotNil(t, mongoClient)
}
