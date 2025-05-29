package messagepersistence

import (
	"gorm.io/gorm"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"
)

type PostgresMessagePersistenceDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

func NewPostgresMessagePersistenceDBContext(
	db *gorm.DB,
) *PostgresMessagePersistenceDBContext {
	// initialize base GormContext
	c := &PostgresMessagePersistenceDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
