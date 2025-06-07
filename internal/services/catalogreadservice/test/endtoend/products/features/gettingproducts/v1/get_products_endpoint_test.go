//go:build e2e
// +build e2e

package v1

import (
	"context"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	httpexpect "github.com/gavv/httpexpect/v2"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestGetAllProducts(t *testing.T) {
	e2eFixture := integration.NewOrderIntegrationTestSharedFixture(t)

	Convey("Get All Products Feature", t, func() {
		e2eFixture.SetupTest()
		ctx := context.Background()

		Convey("Get all products returns ok status", func() {
			Convey("When a request is made to get all products", func() {
				expect := httpexpect.New(t, e2eFixture.BaseAddress)

				Convey("Then the response status should be OK", func() {
					expect.GET("products").
						WithContext(ctx).
						Expect().
						Status(http.StatusOK)
				})
			})
		})

		e2eFixture.TearDownTest()
	})
}
