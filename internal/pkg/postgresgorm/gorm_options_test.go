//go:build unit
// +build unit

// Package postgresgorm provides the gorm options.
package postgresgorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestOptionsName tests the options name.
func TestOptionsName(t *testing.T) {
	assert.Equal(t, "gormOptions", optionName)
}
