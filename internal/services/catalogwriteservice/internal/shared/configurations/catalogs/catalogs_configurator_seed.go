package catalogs

import (
	"time"

	"emperror.dev/errors"
	"gorm.io/gorm"

	gofakeit "github.com/brianvoe/gofakeit/v6"
	uuid "github.com/satori/go.uuid"

	datamodel "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/products/data/datamodels"
)

func (ic *CatalogWriteServiceConfigurator) seedCatalogs(
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
			ID:          uuid.NewV4(),
			Name:        gofakeit.Name(),
			CreatedAt:   time.Now(),
			Description: gofakeit.AdjectiveDescriptive(),
			Price:       gofakeit.Price(100, 1000),
		},
		{
			ID:          uuid.NewV4(),
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
