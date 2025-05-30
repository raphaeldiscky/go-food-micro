// Package postgressqlx provides a set of functions for the postgressqlx.
package postgressqlx

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// optionName is the name of the option.
var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[PostgresSqlxOptions]())

// PostgresSqlxOptions is a struct that contains the postgressqlx options.
type PostgresSqlxOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
}

// provideConfig provides the postgressqlx options.
func provideConfig(environment environment.Environment) (*PostgresSqlxOptions, error) {
	return config.BindConfigKey[*PostgresSqlxOptions](optionName, environment)
}
