package test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/fxapp/contracts"
	config3 "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/mongodb"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/bus"
	config2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/rabbitmq/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
	mongo2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/mongo"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/rabbitmq"
	redis2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/containers/testcontainer/redis"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/config"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/contracts/data"
	catalogs2 "github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/configurations/catalogs"

	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/otel/trace"
)

type TestApp struct{}

type TestAppResult struct {
	Cfg                    *config.Config
	Bus                    bus.RabbitmqBus
	Container              contracts.Container
	Logger                 logger.Logger
	RabbitmqOptions        *config2.RabbitmqOptions
	EchoHTTPOptions        *config3.EchoHTTPOptions
	MongoDbOptions         *mongodb.MongoDbOptions
	RedisOptions           *redis.RedisOptions
	ProductCacheRepository data.ProductCacheRepository
	ProductRepository      data.ProductRepository
	MongoClient            *mongo.Client
	Tracer                 trace.Tracer
}

func NewTestApp() *TestApp {
	return &TestApp{}
}

func (a *TestApp) Run(t *testing.T) (result *TestAppResult) {
	lifetimeCtx := context.Background()

	// ref: https://github.com/uber-go/fx/blob/master/app_test.go
	appBuilder := NewCatalogsReadTestApplicationBuilder(t)
	appBuilder.ProvideModule(catalogs2.CatalogsServiceModule)

	// replace real options with docker container options for testing
	appBuilder.Decorate(rabbitmq.RabbitmqContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(mongo2.MongoContainerOptionsDecorator(t, lifetimeCtx))
	appBuilder.Decorate(redis2.RedisContainerOptionsDecorator(t, lifetimeCtx))

	testApp := appBuilder.Build()

	testApp.ConfigureCatalogs()

	testApp.MapCatalogsEndpoints()

	testApp.ResolveFunc(
		func(cfg *config.Config,
			bus bus.RabbitmqBus,
			logger logger.Logger,
			rabbitmqOptions *config2.RabbitmqOptions,
			mongoOptions *mongodb.MongoDbOptions,
			redisOptions *redis.RedisOptions,
			productCacheRepository data.ProductCacheRepository,
			productRepository data.ProductRepository,
			echoOptions *config3.EchoHTTPOptions,
			mongoClient *mongo.Client,
			tracer trace.Tracer,
		) {
			result = &TestAppResult{
				Bus:                    bus,
				Cfg:                    cfg,
				Container:              testApp,
				Logger:                 logger,
				RabbitmqOptions:        rabbitmqOptions,
				MongoDbOptions:         mongoOptions,
				ProductRepository:      productRepository,
				ProductCacheRepository: productCacheRepository,
				EchoHTTPOptions:        echoOptions,
				MongoClient:            mongoClient,
				RedisOptions:           redisOptions,
				Tracer:                 tracer,
			}
		},
	)

	// we need a longer timout for up and running our testcontainers
	duration := time.Second * 300

	// short timeout for handling start hooks and setup dependencies
	startCtx, cancel := context.WithTimeout(context.Background(), duration)
	defer cancel()
	err := testApp.Start(startCtx)
	if err != nil {
		t.Errorf("Error starting, err: %v", err)
		os.Exit(1)
	}

	t.Cleanup(func() {
		// short timeout for handling stop hooks
		stopCtx, cancel := context.WithTimeout(context.Background(), duration)
		defer cancel()

		err = testApp.Stop(stopCtx)
		require.NoError(t, err)
	})

	return
}
