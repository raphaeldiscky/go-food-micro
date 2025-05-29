package eventstoredb

import (
	"context"
	"testing"

	"github.com/EventStore/EventStore-Client-Go/esdb"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/eventstroredb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
)

func Test_Custom_EventStoreDB_Container(t *testing.T) {
	var esdbClient *esdb.Client
	ctx := context.Background()

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		eventstroredb.ModuleFunc(func() {
		}),
		fx.Decorate(EventstoreDBContainerOptionsDecorator(t, ctx)),
		fx.Populate(&esdbClient),
	).RequireStart()
}
