package postgres

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[PostgresPgxOptions]())

type PostgresPgxOptions struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	User     string `mapstructure:"user"`
	DBName   string `mapstructure:"dbName"`
	SSLMode  bool   `mapstructure:"sslMode"`
	Password string `mapstructure:"password"`
	LogLevel int    `mapstructure:"logLevel"`
}

func provideConfig(environment environment.Environment) (*PostgresPgxOptions, error) {
	return config.BindConfigKey[*PostgresPgxOptions](optionName, environment)
}
