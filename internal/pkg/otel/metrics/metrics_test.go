//go:build integration
// +build integration

// Package metrics provides a test for the metrics.
package metrics_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
)

// TestMetricsEndpoint tests the metrics endpoint.
func TestMetricsEndpoint(t *testing.T) {
	var cfg *metrics.MetricsOptions
	var echoServer contracts.EchoHTTPServer

	app := fxtest.New(
		t,
		config.ModuleFunc(environment.Test), // Use test environment with config.test.json
		zap.Module,
		fxlog.FxLogger,
		customEcho.Module,
		metrics.Module,
		fx.Populate(&cfg, &echoServer),
	)
	app.RequireStart()
	defer app.RequireStop()

	// Metrics endpoint is served on the Echo server, not on a separate port
	var metricsPath string
	if cfg.MetricsRoutePath == "" {
		metricsPath = "metrics"
	} else {
		metricsPath = cfg.MetricsRoutePath
	}

	// Get Echo server configuration
	echoOptions := echoServer.Cfg()
	url := fmt.Sprintf("%s%s/%s", echoOptions.Host, echoOptions.Port, metricsPath)

	t.Run("returns status OK", func(t *testing.T) {
		//nolint:gosec // G107: Potential HTTP request made with variable url
		res, err := http.Get(url)
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusOK, res.StatusCode)
	})

	t.Run("returns metrics data", func(t *testing.T) {
		//nolint:gosec // G107: Potential HTTP request made with variable url
		res, err := http.Get(url)
		require.NoError(t, err)
		defer res.Body.Close()

		b, err := io.ReadAll(res.Body)
		require.NoError(t, err)

		// Check that the metrics endpoint returns some prometheus metrics
		require.Contains(t, string(b), "promhttp_metric_handler_requests_total")
	})
}
