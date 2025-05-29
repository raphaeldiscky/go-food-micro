//go:build integration
// +build integration

package queries

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	mediatr "github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/get_product_by_id/v1/queries"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestGetProductById(t *testing.T) {
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(t)
	ctx := context.Background()

	Convey("Getting Product By ID Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.SetupTest(t)

		knownProductID, err := uuid.FromString(integrationTestSharedFixture.Items[0].ID)
		unknownProductID := uuid.NewV4()
		So(err, ShouldBeNil)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey(
			"Returning an existing product with valid ID from the database with correct properties",
			func() {
				Convey("Given a product with a known ID exists in the database", func() {
					query, err := queries.NewGetProductByID(knownProductID)
					So(err, ShouldBeNil)

					Convey(
						"When we execute GetProductByID query for a product with known ID",
						func() {
							result, err := mediatr.Send[*queries.GetProductByID, *dtos.GetProductByIDResponseDto](
								ctx,
								query,
							)

							Convey("Then it should retrieve product successfully", func() {
								So(result, ShouldNotBeNil)
								So(result.Product, ShouldNotBeNil)
								So(err, ShouldBeNil)

								Convey(
									"And the retrieved product should have the correct ID",
									func() {
										// Assert that the retrieved product's ID matches the known ID.
										So(result.Product.ID, ShouldEqual, knownProductID.String())
									},
								)

								Convey("And other product properties should be correct", func() {
									// Assert other properties of the retrieved product as needed.
								})
							})
						},
					)
				})
			},
		)

		Convey("Returning a NotFound error when product with specific id does not exist", func() {
			Convey("Given a product with a unknown ID in the database", func() {
				// Create a test context and an unknown product ID.

				query, err := queries.NewGetProductByID(unknownProductID)
				So(err, ShouldBeNil)

				Convey("When GetProductByID executed for a product with an unknown ID", func() {
					result, err := mediatr.Send[*queries.GetProductByID, *dtos.GetProductByIDResponseDto](
						ctx,
						query,
					)

					Convey("Then the product should not be found and null result", func() {
						// Assert that the error indicates that the product was not found.
						So(result, ShouldBeNil)
						So(err, ShouldNotBeNil)
					})
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})
}
