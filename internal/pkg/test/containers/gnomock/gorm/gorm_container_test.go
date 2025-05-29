// Package gorm provides a gorm container.
package gorm

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGormContainer tests the gorm container.
func TestGormContainer(t *testing.T) {
	gorm, err := NewGnoMockGormContainer().PopulateContainerOptions(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, gorm)
}
