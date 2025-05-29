// Package postgresgorm provides a set of functions for the postgres gorm.
package postgresgorm

import (
	"fmt"
	"path/filepath"

	"github.com/iancoleman/strcase"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	typeMapper "github.com/raphaeldiscky/go-food-micro/internal/pkg/reflection/typemapper"
)

// optionName is the name of the option.
var optionName = strcase.ToLowerCamel(typeMapper.GetGenericTypeNameByT[GormOptions]())

// GormOptions is a struct that contains the gorm options.
type GormOptions struct {
	UseInMemory   bool   `mapstructure:"useInMemory"`
	UseSQLLite    bool   `mapstructure:"useSqlLite"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	User          string `mapstructure:"user"`
	DBName        string `mapstructure:"dbName"`
	SSLMode       bool   `mapstructure:"sslMode"`
	Password      string `mapstructure:"password"`
	EnableTracing bool   `mapstructure:"enableTracing" default:"true"`
}

// Dns returns the dns.
func (h *GormOptions) Dns() string {
	if h.UseInMemory {
		return ""
	}

	if h.UseSQLLite {
		projectRootDir := environment.GetProjectRootWorkingDirectory()
		dbFilePath := filepath.Join(projectRootDir, fmt.Sprintf("%s.db", h.DBName))

		return dbFilePath
	}

	datasource := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		h.User,
		h.Password,
		h.Host,
		h.Port,
		h.DBName,
	)

	return datasource
}

// provideConfig provides the gorm options.
func provideConfig(environment environment.Environment) (*GormOptions, error) {
	return config.BindConfigKey[*GormOptions](optionName, environment)
}
