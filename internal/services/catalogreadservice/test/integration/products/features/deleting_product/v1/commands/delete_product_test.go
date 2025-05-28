//go:build integration
// +build integration

package commands

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	mediatr "github.com/mehdihadeli/go-mediatr"
	uuid "github.com/satori/go.uuid"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/products/features/deleting_products/v1/commands"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/testfixture/integration"
)

func TestDeleteProduct(t *testing.T) {
	integrationTestSharedFixture := integration.NewIntegrationTestSharedFixture(
		t,
	)

	Convey("Deleting Product Feature", t, func() {
		ctx := context.Background()
		integrationTestSharedFixture.SetupTest(t)

		// https://specflow.org/learn/gherkin/#learn-gherkin
		// scenario
		Convey("Deleting an existing product from the database", func() {
			Convey("Given an existing product in the mongo database", func() {
				productId, err := uuid.FromString(
					integrationTestSharedFixture.Items[0].ProductID,
				)
				So(err, ShouldBeNil)

				command, err := commands.NewDeleteProduct(productId)
				So(err, ShouldBeNil)

				Convey("When we execute the DeleteProduct command", func() {
					result, err := mediatr.Send[*commands.DeleteProduct, *mediatr.Unit](
						context.Background(),
						command,
					)

					Convey(
						"Then the product should be deleted successfully in mongo database",
						func() {
							So(err, ShouldBeNil)
							So(result, ShouldNotBeNil)

							Convey(
								"And the product should no longer exist in the system",
								func() {
									deletedProduct, _ := integrationTestSharedFixture.ProductRepository.GetProductByProductID(
										ctx,
										productId.String(),
									)
									So(deletedProduct, ShouldBeNil)
								},
							)
						},
					)
				})
			})
		})

		integrationTestSharedFixture.TearDownTest()
	})
}
