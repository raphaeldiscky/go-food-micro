// Package redis provides a set of functions for the redis package.
package redis

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// optionName is the name of the option for the redis client.
var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[RedisOptions]())

// RedisOptions is a struct that contains the options for the redis client.
type RedisOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Password      string `mapstructure:"password"`
	Database      int    `mapstructure:"database"`
	PoolSize      int    `mapstructure:"poolSize"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

// provideConfig provides the config for the redis client.
func provideConfig(environment environment.Environment) (*RedisOptions, error) {
	return config.BindConfigKey[*RedisOptions](optionName, environment)
}
