// Package contracts provides a health check contracts.
package contracts

import (
	"go.uber.org/fx"
)

// HealthParams is a struct that represents a health params.
type HealthParams struct {
	fx.In

	Healths []Health `group:"healths"`
}
