// Package postgres provides a module for the postgres.
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"emperror.dev/errors"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"

	bun2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/bun"
)

// NewBunDB creates a new bun db.
func NewBunDB(cfg *bun2.BunConfig) (*bun.DB, error) {
	if cfg.DBName == "" {
		return nil, errors.New("dbname is required in the config.")
	}

	err := createDB(cfg)
	if err != nil {
		return nil, err
	}

	// https://bun.uptrace.dev/postgres/#pgdriver
	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.DBName,
	)
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(datasource)))

	// pgconn := pgdriver.NewConnector(
	//	pgdriver.WithNetwork("tcp"),
	//	pgdriver.WithAddr("localhost:5437"),
	//	pgdriver.WithTLSConfig(&tls.config{InsecureSkipVerify: true}),
	//	pgdriver.WithUser("test"),
	//	pgdriver.WithPassword("test"),
	//	pgdriver.WithDatabase("test"),
	//	pgdriver.WithApplicationName("myapp"),
	//	pgdriver.WithTimeout(5*time.Second),
	//	pgdriver.WithDialTimeout(5*time.Second),
	//	pgdriver.WithReadTimeout(5*time.Second),
	//	pgdriver.WithWriteTimeout(5*time.Second),
	//	pgdriver.WithConnParams(map[string]interface{}{
	//		"search_path": "my_search_path",
	//	}),
	//)
	// sqldb := sql.OpenDB(pgconn)

	db := bun.NewDB(sqldb, pgdialect.New())

	return db, nil
}

// createDB creates a new database.
func createDB(cfg *bun2.BunConfig) error {
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
	rows, err := sqldb.QueryContext(
		context.Background(),
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
			log.Printf("Error closing database: %v", err)
		}
	}()

	return nil
}
