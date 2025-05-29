// Package postgresgorm provides a set of functions for the postgres gorm.
package postgresgorm

import (
	"database/sql"
	"fmt"

	"emperror.dev/errors"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/plugin/opentelemetry/tracing"

	gormPostgres "gorm.io/driver/postgres"

	defaultlogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/gromlog"
)

// NewGorm creates a new gorm.
func NewGorm(cfg *GormOptions) (*gorm.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("DBName is required in the config.")
	}

	if cfg.UseSQLLite {
		db, err := createSQLLiteDB(cfg.DND())

		return db, err
	}

	// InMemory doesn't work correctly with transactions - seems when we `Begin` a transaction on gorm.DB (with SQLLite in-memory) our previous gormDB before transaction will remove and the new gormDB with tx will go on the memory
	if cfg.UseInMemory {
		db, err := createInMemoryDB()

		return db, err
	}

	err := createPostgresDB(cfg)
	if err != nil {
		return nil, err
	}

	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	gormDb, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		},
	)
	if err != nil {
		return nil, err
	}

	// add tracing to gorm
	if cfg.EnableTracing {
		err = gormDb.Use(tracing.NewPlugin())
		if err != nil {
			return nil, err
		}
	}

	return gormDb, nil
}

// createInMemoryDB creates a new in memory db.
func createInMemoryDB() (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	db, err := gorm.Open(
		sqlite.Open(":memory:"),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		})

	return db, err
}

// createSQLLiteDB creates a new sqlite db.
func createSQLLiteDB(dbFilePath string) (*gorm.DB, error) {
	// https://gorm.io/docs/connecting_to_the_database.html#SQLite
	// https://github.com/glebarez/sqlite
	// https://www.connectionstrings.com/sqlite/
	gormSQLLiteDB, err := gorm.Open(
		sqlite.Open(dbFilePath),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		})

	return gormSQLLiteDB, err
}

// NewSQLDB creates a new sql db.
func NewSQLDB(orm *gorm.DB) (*sql.DB, error) {
	return orm.DB()
}

// createPostgresDB creates a new postgres db.
func createPostgresDB(cfg *GormOptions) error {
	var db *sql.DB

	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	dataSourceName := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		"postgres",
		cfg.Password,
	)
	postgresGormDB, err := gorm.Open(
		gormPostgres.Open(dataSourceName),
		&gorm.Config{
			Logger: gromlog.NewGormCustomLogger(defaultlogger.GetLogger()),
		},
	)
	if err != nil {
		return err
	}

	db, err = postgresGormDB.DB()
	if err != nil {
		return err
	}

	rows, err := db.Query(
		fmt.Sprintf(
			"SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'",
			cfg.DBName,
		),
	)
	if err != nil {
		return err
	}

	var exists int
	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = db.Exec(fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer db.Close()

	return nil
}
