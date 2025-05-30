// Package metrics provides a test for the metrics.
package metrics_test

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	"go.uber.org/fx"
	"go.uber.org/fx/fxtest"

	ginkgo "github.com/onsi/ginkgo/v2"
	gomega "github.com/onsi/gomega"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	customEcho "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/customecho"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/otel/metrics"
)

// TestHealth tests the health of the metrics.
func TestHealth(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)

	ginkgo.RunSpecs(t, "/health suite")
}

var _ = ginkgo.Describe("/", ginkgo.Ordered, func() {
	var (
		url string
		err error
		res *http.Response
	)

	ginkgo.BeforeAll(func() {
		var cfg *metrics.MetricsOptions

		fxtest.New(
			ginkgo.GinkgoT(),
			zap.Module,
			fxlog.FxLogger,
			config.Module,
			customEcho.Module,

			metrics.Module,

			fx.Populate(&cfg),
		).RequireStart()

		url = fmt.Sprintf("http://%s:%s/metrics", cfg.Host, cfg.Port)
	})

	ginkgo.BeforeEach(func() {
		//nolint:gosec // G107: Potential HTTP request made with variable url
		res, err = http.Get(url)
	})
	ginkgo.It("returns status OK", func() {
		gomega.Expect(err).To(gomega.BeNil())
		gomega.Expect(res.StatusCode).To(gomega.Equal(http.StatusOK))
	})

	ginkgo.It("returns how many requests were made", func() {
		b, err := io.ReadAll(res.Body)
		gomega.Expect(err).To(gomega.BeNil())

		gomega.Expect(
			b,
		).To(gomega.ContainSubstring(`promhttp_metric_handler_requests_total{code="200"} 1`))
	})
})
