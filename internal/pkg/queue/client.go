// Package queue provides a set of functions for the queue.
package queue

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"go.uber.org/fx"

	redis2 "github.com/raphaeldiscky/go-food-micro/internal/pkg/redis"
)

// NewClient creates a new client.
func NewClient(config *redis2.RedisOptions) *asynq.Client {
	return asynq.NewClient(
		asynq.RedisClientOpt{Addr: fmt.Sprintf("%s:%d", config.Host, config.Port)},
	)
}

// HookClient hooks the client.
func HookClient(lifecycle fx.Lifecycle, client *asynq.Client) {
	lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			return client.Close()
		},
	})
}
