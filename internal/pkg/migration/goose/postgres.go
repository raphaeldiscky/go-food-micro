// Package goose provides a migration runner.
package goose

import (
	"context"
	"database/sql"
	"errors"
	"strconv"

	goose "github.com/pressly/goose/v3"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	migration "github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
)

// goosePostgresMigrator is a goose postgres migrator.
type goosePostgresMigrator struct {
	config *migration.MigrationOptions
	db     *sql.DB
	logger logger.Logger
}

// NewGoosePostgres creates a new goose postgres migrator.
func NewGoosePostgres(
	config *migration.MigrationOptions,
	db *sql.DB,
	logger logger.Logger,
) contracts.PostgresMigrationRunner {
	goose.SetBaseFS(nil)

	return &goosePostgresMigrator{config: config, db: db, logger: logger}
}

// Up runs the up migration.
func (m *goosePostgresMigrator) Up(_ context.Context, version uint) error {
	err := m.executeCommand(migration.Up, version)

	return err
}

// Down runs the down migration.
func (m *goosePostgresMigrator) Down(_ context.Context, version uint) error {
	err := m.executeCommand(migration.Down, version)

	return err
}

// executeCommand executes a command.
func (m *goosePostgresMigrator) executeCommand(
	command migration.CommandType,
	version uint,
) error {
	switch command {
	case migration.Up:
		if version == 0 {
			// In test environment, we need a fix for applying application working directory correctly. we will apply this in our environment setup process in `config/environment` file
			return goose.Run("up", m.db, m.config.MigrationsDir)
		}

		return goose.Run(
			"up-to VERSION ",
			m.db,
			m.config.MigrationsDir,
			strconv.FormatUint(uint64(version), 10),
		)
	case migration.Down:
		if version == 0 {
			return goose.Run("down", m.db, m.config.MigrationsDir)
		}

		return goose.Run(
			"down-to VERSION ",
			m.db,
			m.config.MigrationsDir,
			strconv.FormatUint(uint64(version), 10),
		)
	default:
		return errors.New("invalid migration direction")
	}
}
