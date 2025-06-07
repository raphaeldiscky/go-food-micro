//go:build integration
// +build integration

package v1

import (
	"context"
	"fmt"
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	mediatr "github.com/mehdihadeli/go-mediatr"
	customErrors "github.com/raphaeldiscky/go-food-micro/internal/pkg/http/httperrors/customerrors"
	uuid "github.com/satori/go.uuid"

	getProductByIdQuery "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/features/gettingproductbyid/v1/dtos"
	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/testfixtures/integration"
)

var integrationFixture *integration.CatalogWriteIntegrationTestSharedFixture

func TestGetProductById(t *testing.T) {
	RegisterFailHandler(Fail)
	integrationFixture = integration.NewCatalogWriteIntegrationTestSharedFixture(t)
	RunSpecs(t, "Get Product By ID Integration Tests")
}

var _ = Describe("Get Product by ID Feature", func() {
	var (
		ctx    context.Context
		id     uuid.UUID
		query  *getProductByIdQuery.GetProductByID
		result *dtos.GetProductByIDResponseDto
		err    error
	)

	_ = BeforeEach(func() {
		By("Seeding the required data")
		// call base SetupTest hook before running child hook
		integrationFixture.SetupTest()

		// child hook codes should be here
		id = integrationFixture.Items[0].ID
	})

	_ = AfterEach(func() {
		By("Cleanup test data")
		// call base TearDownTest hook before running child hook
		integrationFixture.TearDownTest()

		// child hook codes should be here
	})

	_ = BeforeSuite(func() {
		ctx = context.Background()

		// in test mode we set rabbitmq `AutoStart=false` in configuration in rabbitmqOptions, so we should run rabbitmq bus manually
		err = integrationFixture.Bus.Start(context.Background())
		Expect(err).ShouldNot(HaveOccurred())

		// wait for consumers ready to consume before publishing messages, preparation background workers takes a bit time (for preventing messages lost)
		time.Sleep(1 * time.Second)
	})

	_ = AfterSuite(func() {
		integrationFixture.Log.Info("TearDownSuite started")
		err := integrationFixture.Bus.Stop()
		Expect(err).ShouldNot(HaveOccurred())
		time.Sleep(1 * time.Second)
	})

	// "Scenario" step for testing returning an existing product with correct properties
	Describe(
		"Returning an existing product with valid ID from the database with correct properties",
		func() {
			Context("Given products exists in the database", func() {
				BeforeEach(func() {
					query, err = getProductByIdQuery.NewGetProductByIDWithValidation(id)
				})

				// "When" step
				When(
					"the GteProductById query is executed for existing product",
					func() {
						BeforeEach(func() {
							result, err = mediatr.Send[*getProductByIdQuery.GetProductByID, *dtos.GetProductByIDResponseDto](
								ctx,
								query,
							)
						})

						// "Then" step
						It("Should not return an error", func() {
							Expect(err).To(BeNil())
						})

						It("Should return a non-nil result", func() {
							Expect(result).NotTo(BeNil())
						})

						It(
							"Should return a product with the correct ID",
							func() {
								Expect(result.Product).NotTo(BeNil())
								Expect(result.Product.ID).To(Equal(id))
							},
						)
					},
				)
			})
		},
	)

	// "Scenario" step for testing returning a NotFound error when record does not exist
	Describe(
		"Returning a NotFound error when product with specific id does not exist",
		func() {
			Context("Given products does not exists in the database", func() {
				BeforeEach(func() {
					// Generate a random UUID that does not exist in the database
					id = uuid.NewV4()
					query, err = getProductByIdQuery.NewGetProductByIDWithValidation(id)
					Expect(err).To(BeNil())
				})

				// "When" step
				When(
					"the GetProductByID query is executed for non-existing product",
					func() {
						BeforeEach(func() {
							result, err = mediatr.Send[*getProductByIdQuery.GetProductByID, *dtos.GetProductByIDResponseDto](
								ctx,
								query,
							)
						})

						// "Then" step
						It("Should return an error", func() {
							Expect(err).To(HaveOccurred())
						})

						It("Should return a NotFound error", func() {
							fmt.Println(err)
							Expect(
								err,
							).To(MatchError(ContainSubstring(fmt.Sprintf("product with id `%s` not found in the database", query.ProductID.String()))))
						})

						It("Should return a custom NotFound error", func() {
							Expect(
								customErrors.IsNotFoundError(err),
							).To(BeTrue())
						})

						It("Should not return a result", func() {
							Expect(result).To(BeNil())
						})
					},
				)
			})
		},
	)
})
