// Package consumer provides a module for the consumer.
package consumer

// ConsumerHandlerConfigurationBuilderFunc is a consumer handler configuration builder function.
type ConsumerHandlerConfigurationBuilderFunc func(ConsumerHandlerConfigurationBuilder)

// ConsumerHandlerConfigurationBuilder is a consumer handler configuration builder.
type ConsumerHandlerConfigurationBuilder interface {
	AddHandler(handler ConsumerHandler) ConsumerHandlerConfigurationBuilder
	Build() *ConsumerHandlersConfiguration
}

// consumerHandlerConfigurationBuilder is a consumer handler configuration builder.
type consumerHandlerConfigurationBuilder struct {
	consumerHandlersConfiguration *ConsumerHandlersConfiguration
}

// NewConsumerHandlersConfigurationBuilder creates a new consumer handlers configuration builder.
func NewConsumerHandlersConfigurationBuilder() ConsumerHandlerConfigurationBuilder {
	return &consumerHandlerConfigurationBuilder{
		consumerHandlersConfiguration: &ConsumerHandlersConfiguration{},
	}
}

// AddHandler adds a handler to the consumer handlers configuration builder.
func (c *consumerHandlerConfigurationBuilder) AddHandler(
	handler ConsumerHandler,
) ConsumerHandlerConfigurationBuilder {
	c.consumerHandlersConfiguration.Handlers = append(
		c.consumerHandlersConfiguration.Handlers,
		handler,
	)

	return c
}

// Build builds the consumer handlers configuration.
func (c *consumerHandlerConfigurationBuilder) Build() *ConsumerHandlersConfiguration {
	return c.consumerHandlersConfiguration
}
