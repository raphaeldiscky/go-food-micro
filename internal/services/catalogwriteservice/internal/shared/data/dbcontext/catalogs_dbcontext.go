package dbcontext

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

type CatalogsGormDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

func NewCatalogsDBContext(db *gorm.DB) *CatalogsGormDBContext {
	// initialize base GormContext
	c := &CatalogsGormDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
