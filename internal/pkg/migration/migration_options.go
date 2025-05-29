// Package migration provides a migration runner.
package migration

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// CommandType is a command type.
type CommandType string

// CommandType constants.
const (
	Up   CommandType = "up"
	Down CommandType = "down"
)

// MigrationOptions is a migration options.
type MigrationOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	VersionTable  string `mapstructure:"versionTable"`
	MigrationsDir string `mapstructure:"migrationsDir"`
	SkipMigration bool   `mapstructure:"skipMigration"`
}

// ProvideConfig provides a migration config.
func ProvideConfig(environment environment.Environment) (*MigrationOptions, error) {
	optionName := strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[MigrationOptions]())

	return config.BindConfigKey[*MigrationOptions](optionName, environment)
}
