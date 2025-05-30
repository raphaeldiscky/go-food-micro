// Package messagepersistence provides a set of functions for the message persistence.
package messagepersistence

import (
	"gorm.io/gorm"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
)

// PostgresMessagePersistenceDBContext is a struct that contains the postgres message persistence db context.
type PostgresMessagePersistenceDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

// NewPostgresMessagePersistenceDBContext creates a new postgres message persistence db context.
func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
) *PostgresMessagePersistenceDBContext {
	// initialize base GormContext
	c := &PostgresMessagePersistenceDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
