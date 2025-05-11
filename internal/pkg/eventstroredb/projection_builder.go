package eventstroredb

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/es/contracts/projection"
)

type projectionsBuilder struct {
	projectionConfiguration *ProjectionsConfigurations
}

func NewProjectionsBuilder() ProjectionsBuilder {
	return &projectionsBuilder{
		projectionConfiguration: &ProjectionsConfigurations{},
	}
}

func (p *projectionsBuilder) AddProjection(projection projection.IProjection) ProjectionsBuilder {
	p.projectionConfiguration.Projections = append(p.projectionConfiguration.Projections, projection)
	return p
}

func (p *projectionsBuilder) AddProjections(projections []projection.IProjection) ProjectionsBuilder {
	p.projectionConfiguration.Projections = append(p.projectionConfiguration.Projections, projections...)
	return p
}

func (p *projectionsBuilder) Build() *ProjectionsConfigurations {
	return p.projectionConfiguration
}
