//go:build integration
// +build integration

package gorm

import (
	"context"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"
	"gorm.io/gorm"

	testcontainers "github.com/testcontainers/testcontainers-go"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
)

// TestCustomGormContainer tests the custom gorm container.
func TestCustomGormContainer(t *testing.T) {
	ctx := context.Background()

	var gorm *gorm.DB

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		gormPostgres.Module,
		fx.Decorate(GormContainerOptionsDecorator(t, ctx)),
		fx.Populate(&gorm),
	).RequireStart()

	assert.NotNil(t, gorm)
}

// TestBuiltinPostgresContainer tests the builtin postgres container.
func TestBuiltinPostgresContainer(t *testing.T) {
	ctx := context.Background()

	// https://github.com/testcontainers/testcontainers-go/blob/f87445303764342cb09ae3cc0e1f80c082b003a4/modules/postgres/postgres_test.go
	ct, err := postgres.RunContainer(
		context.Background(),
		testcontainers.WithImage("postgres"),
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second),
		),
	)
	if err != nil {
		t.Fatal(err)
	}

	// Clean up the container after the test is complete
	t.Cleanup(func() {
		if err := ct.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	})

	host, err := ct.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}
	port, err := ct.MappedPort(ctx, nat.Port("5432/tcp"))
	if err != nil {
		t.Fatal(err)
	}
	gormOptions := &gormPostgres.GormOptions{
		Port:     port.Int(),
		Host:     host,
		Password: "postgres",
		DBName:   "postgres",
		SSLMode:  false,
		User:     "postgres",
	}
	db, err := gormPostgres.NewGorm(gormOptions)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, db)
}
