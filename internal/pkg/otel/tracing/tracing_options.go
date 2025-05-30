// Package tracing provides options for the tracing.
package tracing

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

// TracingOptions is a config for the tracing.
type TracingOptions struct {
	Enabled                   bool                   `mapstructure:"enabled"`
	ServiceName               string                 `mapstructure:"serviceName"`
	Version                   string                 `mapstructure:"version"`
	InstrumentationName       string                 `mapstructure:"instrumentationName"`
	ID                        int64                  `mapstructure:"id"`
	AlwaysOnSampler           bool                   `mapstructure:"alwaysOnSampler"`
	ZipkinExporterOptions     *ZipkinExporterOptions `mapstructure:"zipkinExporterOptions"`
	JaegerExporterOptions     *OTLPProvider          `mapstructure:"jaegerExporterOptions"`
	ElasticApmExporterOptions *OTLPProvider          `mapstructure:"elasticApmExporterOptions"`
	UptraceExporterOptions    *OTLPProvider          `mapstructure:"uptraceExporterOptions"`
	SignozExporterOptions     *OTLPProvider          `mapstructure:"signozExporterOptions"`
	TempoExporterOptions      *OTLPProvider          `mapstructure:"tempoExporterOptions"`
	UseStdout                 bool                   `mapstructure:"useStdout"`
	UseOTLP                   bool                   `mapstructure:"useOTLP"`
	OTLPProviders             []OTLPProvider         `mapstructure:"otlpProviders"`
}

// ZipkinExporterOptions is a config for the zipkin exporter.
type ZipkinExporterOptions struct {
	Url string `mapstructure:"url"`
}

// ProvideTracingConfig provides a tracing config.
func ProvideTracingConfig(
	environment environment.Environment,
) (*TracingOptions, error) {
	optionName := strcase.ToLowerCamel(
		typeMapper.GetGenericTypeNameByT[TracingOptions](),
	)

	return config.BindConfigKey[*TracingOptions](optionName, environment)
}
