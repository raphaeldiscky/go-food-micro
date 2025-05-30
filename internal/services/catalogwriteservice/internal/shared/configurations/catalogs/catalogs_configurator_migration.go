package catalogs

import (
	"context"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
)

func (ic *CatalogWriteServiceConfigurator) migrateCatalogs(
	runner contracts.PostgresMigrationRunner,
) error {
	// - for complex migration and ability to back-track to specific migration revision it is better we use `goose`, but if we want to use built-in gorm migration we can also sync gorm with `atlas` integration migration versioning for getting migration history from grom changes
	// - here I used goose for migration, with using cmd/migration file
	// migration with goose
	return migrateGoose(runner)
}

func migrateGoose(
	runner contracts.PostgresMigrationRunner,
) error {
	err := runner.Up(context.Background(), 0)

	return err
}
