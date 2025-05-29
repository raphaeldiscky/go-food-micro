package postgrespxg

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
)

func Test_Custom_PostgresPgx_Container(t *testing.T) {
	gorm, err := NewPostgresPgxContainers(
		defaultLogger.GetLogger(),
	).PopulateContainerOptions(context.Background(), t)
	require.NoError(t, err)

	assert.NotNil(t, gorm)
}
