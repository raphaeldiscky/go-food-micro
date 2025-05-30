// Package binder provides config binding functionality.
package binder

import (
	"fmt"

	"github.com/spf13/viper"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
)

// BindConfigKey binds the config key.
func BindConfigKey[T any](
	configKey string,
	_ ...environment.Environment,
) (T, error) {
	var cfg T
	if len(configKey) == 0 {
		if err := viper.Unmarshal(&cfg); err != nil {
			return *new(T), fmt.Errorf("viper.Unmarshal: %w", err)
		}
	} else {
		if err := viper.UnmarshalKey(configKey, &cfg); err != nil {
			return *new(T), fmt.Errorf("viper.Unmarshal: %w", err)
		}
	}

	return cfg, nil
}
