// Package postgresgorm provides a set of functions for the postgres gorm.
package postgresgorm

import (
	"context"
	"database/sql"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/health/contracts"
)

// gormHealthChecker is a struct that contains the gorm health checker.
type gormHealthChecker struct {
	client *sql.DB
}

// NewGormHealthChecker creates a new gorm health checker.
func NewGormHealthChecker(client *sql.DB) contracts.Health {
	return &gormHealthChecker{client}
}

// CheckHealth checks the health of the gorm.
func (healthChecker *gormHealthChecker) CheckHealth(ctx context.Context) error {
	return healthChecker.client.PingContext(ctx)
}

// GetHealthName returns the health name.
func (healthChecker *gormHealthChecker) GetHealthName() string {
	return "postgres"
}
