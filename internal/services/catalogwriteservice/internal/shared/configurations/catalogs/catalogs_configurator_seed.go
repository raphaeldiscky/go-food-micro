package catalogs

import (
	"time"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"

	"emperror.dev/errors"
	"github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (ic *CatalogsServiceConfigurator) seedCatalogs(
	db *gorm.DB,
) error {
	err := seedDataManually(db)
	if err != nil {
		return err
	}

	return nil
}

func seedDataManually(gormDB *gorm.DB) error {
	var count int64

	// https://gorm.io/docs/advanced_query.html#Count
	gormDB.Model(&datamodel.ProductDataModel{}).Count(&count)
	if count > 0 {
		return nil
	}

	products := []*datamodel.ProductDataModel{
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			Id:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
	}

	err := gormDB.CreateInBatches(products, len(products)).Error
	if err != nil {
		return errors.Wrap(err, "error in seed database")
	}

	return nil
}
