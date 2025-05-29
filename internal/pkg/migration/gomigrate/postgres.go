// Package gomigrate provides a migration runner.
package gomigrate

import (
	"context"
	"database/sql"
	"fmt"

	"emperror.dev/errors"

	migrate "github.com/golang-migrate/migrate/v4"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
)

// goMigratePostgresMigrator is a go migrate postgres migrator.
type goMigratePostgresMigrator struct {
	config     *migration.MigrationOptions
	db         *sql.DB
	datasource string
	logger     logger.Logger
	migration  *migrate.Migrate
}

// NewGoMigratorPostgres creates a new go migrate postgres migrator.
func NewGoMigratorPostgres(
	config *migration.MigrationOptions,
	db *sql.DB,
	logger logger.Logger,
) (contracts.PostgresMigrationRunner, error) {
	if config.DBName == "" {
		return nil, errors.New("dbname is required in the config")
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.DBName,
	)

	// In test environment, ewe need a fix for applying application working directory correctly. we will apply this in our environment setup process in `config/environment` file
	migration, err := migrate.New(fmt.Sprintf("file://%s", config.MigrationsDir), datasource)
	if err != nil {
		return nil, errors.WrapIf(err, "failed to initialize migrator")
	}

	return &goMigratePostgresMigrator{
		config:     config,
		db:         db,
		datasource: datasource,
		logger:     logger,
		migration:  migration,
	}, nil
}

// Up runs the up migration.
func (m *goMigratePostgresMigrator) Up(_ context.Context, version uint) error {
	if m.config.SkipMigration {
		m.logger.Info("database migration skipped")

		return nil
	}

	err := m.executeCommand(migration.Up, version)

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	m.logger.Info("migration finished")

	return nil
}

// Down runs the down migration.
func (m *goMigratePostgresMigrator) Down(_ context.Context, version uint) error {
	if m.config.SkipMigration {
		m.logger.Info("database migration skipped")

		return nil
	}

	err := m.executeCommand(migration.Up, version)

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}

	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	m.logger.Info("migration finished")

	return nil
}

// executeCommand executes a command.
func (m *goMigratePostgresMigrator) executeCommand(
	command migration.CommandType,
	version uint,
) error {
	var err error
	switch command {
	case migration.Up:
		if version == 0 {
			err = m.migration.Up()
		} else {
			err = m.migration.Migrate(version)
		}
	case migration.Down:
		if version == 0 {
			err = m.migration.Down()
		} else {
			err = m.migration.Migrate(version)
		}
	default:
		err = errors.New("invalid migration direction")
	}

	if errors.Is(err, migrate.ErrNoChange) {
		return nil
	}
	if err != nil {
		return errors.WrapIf(err, "failed to migrate database")
	}

	return nil
}
