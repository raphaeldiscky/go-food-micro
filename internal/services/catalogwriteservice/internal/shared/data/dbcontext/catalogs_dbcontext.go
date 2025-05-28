package dbcontext

import (
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm/gormdbcontext"

	"gorm.io/gorm"
)

// CatalogsGormDBContext is a struct that contains the CatalogsGormDBContext
type CatalogsGormDBContext struct {
	// our dbcontext base
	contracts.GormDBContext
}

// NewCatalogsDBContext is a constructor for the CatalogsGormDBContext
func NewCatalogsDBContext(db *gorm.DB) *CatalogsGormDBContext {
	// initialize base GormContext
	c := &CatalogsGormDBContext{GormDBContext: gormdbcontext.NewGormDBContext(db)}

	return c
}
