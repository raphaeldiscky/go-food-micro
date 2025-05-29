// Package contracts provides a set of contracts for the postgresgorm package.
package contracts

import (
	"context"

	"gorm.io/gorm"
)

// GormContext is a context that contains a gorm.DB.
type GormContext struct {
	Tx *gorm.DB
	context.Context
}
