// Package pipeline provides a consumer pipeline configuration builder.
package pipeline

// ConsumerPipelineConfigurationBuilderFunc is a function that builds a consumer pipeline configuration.
type ConsumerPipelineConfigurationBuilderFunc func(ConsumerPipelineConfigurationBuilder)

// ConsumerPipelineConfigurationBuilder is a type that represents a consumer pipeline configuration builder.
type ConsumerPipelineConfigurationBuilder interface {
	AddPipeline(pipeline ConsumerPipeline) ConsumerPipelineConfigurationBuilder
	Build() *ConsumerPipelineConfiguration
}

// consumerPipelineConfigurationBuilder is a struct that represents a consumer pipeline configuration builder.
type consumerPipelineConfigurationBuilder struct {
	pipelineConfigurations *ConsumerPipelineConfiguration
}

// NewConsumerPipelineConfigurationBuilder is a function that creates a new consumer pipeline configuration builder.
func NewConsumerPipelineConfigurationBuilder() ConsumerPipelineConfigurationBuilder {
	return &consumerPipelineConfigurationBuilder{
		pipelineConfigurations: &ConsumerPipelineConfiguration{},
	}
}

// AddPipeline is a function that adds a pipeline to the consumer pipeline configuration.
func (c *consumerPipelineConfigurationBuilder) AddPipeline(
	pipeline ConsumerPipeline,
) ConsumerPipelineConfigurationBuilder {
	c.pipelineConfigurations.Pipelines = append(c.pipelineConfigurations.Pipelines, pipeline)

	return c
}

// Build is a function that builds the consumer pipeline configuration.
func (c *consumerPipelineConfigurationBuilder) Build() *ConsumerPipelineConfiguration {
	return c.pipelineConfigurations
}
