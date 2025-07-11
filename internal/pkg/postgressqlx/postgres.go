// Package postgressqlx provides a set of functions for the postgressqlx.
package postgressqlx

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"emperror.dev/errors"
	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/bun/driver/pgdriver"

	_ "github.com/jackc/pgx/v4/stdlib" // load pgx driver for PostgreSQL

	goqu "github.com/doug-martin/goqu/v9"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
)

// Sqlx is a struct that contains the sqlx.
type Sqlx struct {
	SqlxDB          *sqlx.DB
	DB              *sql.DB
	SquirrelBuilder squirrel.StatementBuilderType
	GoquBuilder     *goqu.SelectDataset
	config          *PostgresSqlxOptions
	logger          logger.Logger
}

// NewSqlxConn creates a database connection with appropriate pool configuration.
// and runs migration to prepare database.
//
// Migration will be omitted if appropriate config parameter set.
func NewSqlxConn(cfg *PostgresSqlxOptions) (*Sqlx, error) {
	// Define database connection settings.
	maxConn, err := strconv.Atoi(os.Getenv("DB_MAX_CONNECTIONS"))
	if err != nil {
		return nil, errors.New("error in converting DB_MAX_CONNECTIONS to int")
	}
	maxIdleConn, err := strconv.Atoi(os.Getenv("DB_MAX_IDLE_CONNECTIONS"))
	if err != nil {
		return nil, errors.New("error in converting DB_MAX_IDLE_CONNECTIONS to int")
	}
	maxLifetimeConn, err := strconv.Atoi(os.Getenv("DB_MAX_LIFETIME_CONNECTIONS"))
	if err != nil {
		return nil, errors.New("error in converting DB_MAX_LIFETIME_CONNECTIONS to int")
	}

	var dataSourceName string

	if cfg.DBName == "" {
		return nil, errors.New("dbname is required in the config.")
	}

	err = createDB(cfg)
	if err != nil {
		return nil, err
	}

	dataSourceName = fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s",
		cfg.Host,
		cfg.Port,
		cfg.User,
		cfg.DBName,
		cfg.Password,
	)

	// Define database connection for PostgreSQL.
	db, err := sqlx.Connect("pgx", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("error, not connected to database, %w", err)
	}

	// stdlib package doesn't have a compat layer for pgxpool
	// so had to use standard sql api for pool configuration.
	db.SetMaxOpenConns(maxConn)                           // the defaultLogger is 0 (unlimited)
	db.SetMaxIdleConns(maxIdleConn)                       // defaultMaxIdleConns = 2
	db.SetConnMaxLifetime(time.Duration(maxLifetimeConn)) // 0, connections are reused forever

	// Try to ping database.
	if err := db.PingContext(context.Background()); err != nil {
		defer func() {
			if err := db.Close(); err != nil {
				log.Fatalf("Error closing database: %v", err)
			}
		}()

		return nil, fmt.Errorf("error, not sent ping to database, %w", err)
	}

	// squirrel
	squirrelBuilder := squirrel.StatementBuilder.
		PlaceholderFormat(squirrel.Dollar).RunWith(db)

	// goqu
	dialect := goqu.Dialect("postgres")
	database := dialect.DB(db)
	goquBuilder := database.From()

	sqlx := &Sqlx{
		DB:              db.DB,
		SqlxDB:          db,
		SquirrelBuilder: squirrelBuilder,
		GoquBuilder:     goquBuilder,
		config:          cfg,
	}

	return sqlx, nil
}

func createDB(cfg *PostgresSqlxOptions) error {
	// we should choose a default database in the connection, but because we don't have a database yet we specify postgres default database 'postgres'
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		"postgres",
	)

	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	var exists int
	rows, err := sqldb.QueryContext(context.Background(),
		fmt.Sprintf("SELECT 1 FROM  pg_catalog.pg_database WHERE datname='%s'", cfg.DBName),
	)
	if err != nil {
		return err
	}

	defer func() {
		if err := rows.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error checking database existence: %w", err)
	}

	if rows.Next() {
		err = rows.Scan(&exists)
		if err != nil {
			return err
		}
	}

	if exists == 1 {
		return nil
	}

	_, err = sqldb.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", cfg.DBName))
	if err != nil {
		return err
	}

	defer func() {
		if err := sqldb.Close(); err != nil {
			log.Fatalf("Error closing database: %v", err)
		}
	}()

	return nil
}

// Close closes the database connection.
func (db *Sqlx) Close() {
	if err := db.DB.Close(); err != nil {
		db.logger.Error(
			"error in closing database: %v",
			err,
		)
	}
	if err := db.SqlxDB.Close(); err != nil {
		db.logger.Error(
			"error in closing sqlx database: %v",
			err,
		)
	}
}

// Ref:https://dev.to/techschoolguru/a-clean-way-to-implement-database-transaction-in-golang-2ba

// ExecTx executes a transaction with provided function.
func (db *Sqlx) ExecTx(ctx context.Context, fn func(*Sqlx) error) error {
	tx, err := db.DB.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelReadCommitted})
	if err != nil {
		return err
	}

	err = fn(db)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %w, rb err: %w", err, rbErr)
		}

		return err
	}

	return tx.Commit()
}
