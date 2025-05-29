// Package mongodb provides options for the mongodb.
package mongodb

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// MongoDbOptions is a mongodb options.
type MongoDbOptions struct {
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	Password      string `mapstructure:"password"`
	Database      string `mapstructure:"database"`
	UseAuth       bool   `mapstructure:"useAuth"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

// provideConfig provides a mongodb config.
func provideConfig(
	environment environment.Environment,
) (*MongoDbOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetGenericTypeNameByT[MongoDbOptions](),
	)

	return config.BindConfigKey[*MongoDbOptions](optionName, environment)
}
