// Package eventstoredb provides a eventstoredb container.
package eventstoredb

import (
	"context"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	kdb "github.com/kurrent-io/KurrentDB-Client-Go/kurrentdb"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstoredb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
)

// TestCustomEventStoreDBContainer tests the custom eventstoredb container.
func TestCustomEventStoreDBContainer(t *testing.T) {
	var esdbClient *kdb.Client
	ctx := context.Background()

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		eventstoredb.ModuleFunc(func() {
		}),
		fx.Decorate(EventstoreDBContainerOptionsDecorator(t, ctx)),
		fx.Populate(&esdbClient),
	).RequireStart()
}
