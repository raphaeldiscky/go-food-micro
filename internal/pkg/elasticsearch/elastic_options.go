// Package elasticsearch provides the elasticsearch options.
package elasticsearch

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// optionName is the name of the elasticsearch options.
var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[ElasticOptions]())

// ElasticOptions is a struct that represents the elasticsearch options.
type ElasticOptions struct {
	URL string `mapstructure:"url"`
}

// provideConfig provides the elasticsearch options.
func provideConfig(environment environment.Environment) (*ElasticOptions, error) {
	return config.BindConfigKey[*ElasticOptions](optionName, environment)
}
