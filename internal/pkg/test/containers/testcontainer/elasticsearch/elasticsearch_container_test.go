//go:build integration
// +build integration

// Package elasticsearch provides an elasticsearch container test.
package elasticsearch

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	elasticsearch "github.com/elastic/go-elasticsearch/v8"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/core"
	elasticsearchPkg "github.com/raphaeldiscky/go-food-micro/internal/pkg/elasticsearch"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
)

// TestCustomElasticsearchContainer tests the custom elasticsearch container.
func TestCustomElasticsearchContainer(t *testing.T) {
	ctx := context.Background()

	var elasticsearchClient *elasticsearch.Client

	fxtest.New(t,
		config.ModuleFunc(environment.Test),
		zap.Module,
		fxlog.FxLogger,
		core.Module,
		elasticsearchPkg.Module,
		fx.Decorate(ElasticsearchContainerOptionsDecorator(t, ctx)),
		fx.Populate(&elasticsearchClient),
	).RequireStart()

	assert.NotNil(t, elasticsearchClient)
}
