// Package metrics provides options for the metrics.
package metrics

import (
	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// OTLPProvider is a provider for the otlp.
type OTLPProvider struct {
	Name         string            `mapstructure:"name"`
	Enabled      bool              `mapstructure:"enabled"`
	OTLPEndpoint string            `mapstructure:"otlpEndpoint"`
	OTLPHeaders  map[string]string `mapstructure:"otlpHeaders"`
}

// MetricsOptions is a config for the metrics.
type MetricsOptions struct {
	Host                      string         `mapstructure:"host"`
	Port                      string         `mapstructure:"port"`
	ServiceName               string         `mapstructure:"serviceName"`
	Version                   string         `mapstructure:"version"`
	MetricsRoutePath          string         `mapstructure:"metricsRoutePath"`
	EnableHostMetrics         bool           `mapstructure:"enableHostMetrics"`
	UseStdout                 bool           `mapstructure:"useStdout"`
	InstrumentationName       string         `mapstructure:"instrumentationName"`
	UseOTLP                   bool           `mapstructure:"useOTLP"`
	OTLPProviders             []OTLPProvider `mapstructure:"otlpProviders"`
	ElasticApmExporterOptions *OTLPProvider  `mapstructure:"elasticApmExporterOptions"`
	UptraceExporterOptions    *OTLPProvider  `mapstructure:"uptraceExporterOptions"`
	SignozExporterOptions     *OTLPProvider  `mapstructure:"signozExporterOptions"`
}

// ProvideMetricsConfig provides a metrics config.
func ProvideMetricsConfig(
	environment environment.Environment,
) (*MetricsOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetGenericTypeNameByT[MetricsOptions](),
	)

	return config.BindConfigKey[*MetricsOptions](optionName, environment)
}
