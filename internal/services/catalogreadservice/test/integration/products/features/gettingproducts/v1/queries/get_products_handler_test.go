//go:build integration
// +build integration

package queries

import (
	"context"
	"testing"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/utils"

	. "github.com/smartystreets/goconvey/convey"

	mediatr "github.com/mehdihadeli/go-mediatr"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/gettingproducts/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/gettingproducts/v1/queries"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestGetProducts(t *testing.T) {
	integrationTestSharedFixture := integration.NewCatalogReadIntegrationTestSharedFixture(t)

	Convey("Get All Products Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.SetupTest(t)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Getting a list of existing products from the database", func() {
			Convey("Given a set of existing products in the system", func() {
				query := queries.NewGetProducts(utils.NewListQuery(10, 1))

				Convey("When GetProduct query executed for existing products", func() {
					queryResult, err := mediatr.Send[*queries.GetProducts, *dtos.GetProductsResponseDto](
						ctx,
						query,
					)

					Convey("Then the products should be retrieved successfully", func() {
						// Assert that the error is nil (indicating success).
						So(err, ShouldBeNil)
						So(queryResult, ShouldNotBeNil)
						So(queryResult.Products, ShouldNotBeNil)

						Convey("And the list of products should not be empty", func() {
							// Assert that the list of products is not empty.
							So(queryResult.Products.Items, ShouldNotBeEmpty)

							Convey("And each product should have the correct properties", func() {
								for _, product := range queryResult.Products.Items {
									// Assert properties of each product as needed.
									// For example:
									So(product.Name, ShouldNotBeEmpty)
									So(product.Price, ShouldBeGreaterThan, 0.0)
									// Add more assertions as needed.
								}
							})
						})
					})
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})
}
