// Package hypothesis provides a hypothesis.
package hypothesis

import (
	"context"
	"time"

	"github.com/stretchr/testify/assert"

	reflect "github.com/goccy/go-reflect"

	testUtils "github.com/raphaeldiscky/go-food-micro/internal/pkg/test/utils"
)

// Hypothesis is a interface that represents a hypothesis.
type Hypothesis[T any] interface {
	Validate(ctx context.Context, message string, time time.Duration)
	Test(ctx context.Context, item T)
}

// hypothesis is a struct that represents a hypothesis.
type hypothesis[T any] struct {
	data      T
	condition func(item T) bool
	t         assert.TestingT
}

// Validate validates the hypothesis.
func (h *hypothesis[T]) Validate(ctx context.Context, message string, time time.Duration) {
	err := testUtils.WaitUntilConditionMet(func() bool {
		return !reflect.ValueOf(h.data).IsZero()
	}, time)
	if err != nil {
		assert.FailNowf(h.t, "hypothesis validation failed, %s", message)
	}
}

// Test tests the hypothesis.
func (h *hypothesis[T]) Test(ctx context.Context, item T) {
	h.data = item
}
