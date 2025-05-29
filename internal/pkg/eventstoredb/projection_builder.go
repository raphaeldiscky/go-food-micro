// Package eventstoredb provides a serializer for EventStoreDB.
package eventstoredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

// projectionsBuilder is a struct that represents a projections builder.
type projectionsBuilder struct {
	projectionConfiguration *ProjectionsConfigurations
}

// NewProjectionsBuilder creates a new projections builder.
func NewProjectionsBuilder() ProjectionsBuilder {
	return &projectionsBuilder{
		projectionConfiguration: &ProjectionsConfigurations{},
	}
}

// AddProjection adds a projection to the projections builder.
func (p *projectionsBuilder) AddProjection(projection projection.IProjection) ProjectionsBuilder {
	p.projectionConfiguration.Projections = append(
		p.projectionConfiguration.Projections,
		projection,
	)

	return p
}

// AddProjections adds a projections to the projections builder.
func (p *projectionsBuilder) AddProjections(
	projections []projection.IProjection,
) ProjectionsBuilder {
	p.projectionConfiguration.Projections = append(
		p.projectionConfiguration.Projections,
		projections...)

	return p
}

// Build builds the projections builder.
func (p *projectionsBuilder) Build() *ProjectionsConfigurations {
	return p.projectionConfiguration
}
