// Package environment provides a module for the environment.
package environment

import (
	"log"
	"os"
	"path/filepath"
	"syscall"

	"emperror.dev/errors"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/constants"
)

// Environment is a type for the environment.
type Environment string

var (
	// Development is the development environment.
	Development = Environment(constants.Dev)
	// Test is the test environment.
	Test = Environment(constants.Test)
	// Production is the production environment.
	Production = Environment(constants.Production)
)

// ConfigAppEnv configures the application environment.
func ConfigAppEnv(environments ...Environment) Environment {
	environment := Environment("")
	if len(environments) > 0 {
		environment = environments[0]
	} else {
		environment = Development
	}

	// setup viper to read from os environment with `viper.Get`
	viper.AutomaticEnv()

	// https://articles.wesionary.team/environment-variable-configuration-in-your-golang-project-using-viper-4e8289ef664d
	// load environment variables form .env files to system environment variables, it just finds `.env` file in our current `executing working directory` in our app for example `catalogs_service`
	err := loadEnvFilesRecursive()
	if err != nil {
		log.Printf(".env file cannot be found, err: %v", err)
	}

	setRootWorkingDirectoryEnvironment()

	FixProjectRootWorkingDirectoryPath()

	manualEnv := os.Getenv(constants.AppEnv)

	if manualEnv != "" {
		environment = Environment(manualEnv)
	}

	return environment
}

// IsDevelopment checks if the environment is development.
func (env Environment) IsDevelopment() bool {
	return env == Development
}

// IsProduction checks if the environment is production.
func (env Environment) IsProduction() bool {
	return env == Production
}

// IsTest checks if the environment is test.
func (env Environment) IsTest() bool {
	return env == Test
}

// GetEnvironmentName gets the environment name.
func (env Environment) GetEnvironmentName() string {
	return string(env)
}

// EnvString gets the environment string.
func EnvString(key, fallback string) string {
	if value, ok := syscall.Getenv(key); ok {
		return value
	}

	return fallback
}

// loadEnvFilesRecursive loads the environment files recursively.
func loadEnvFilesRecursive() error {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	// Keep traversing up the directory hierarchy until you find an ".env" file
	for {
		envFilePath := filepath.Join(dir, ".env")
		err := godotenv.Load(envFilePath)

		if err == nil {
			// .env file found and loaded
			return nil
		}

		// Move up one directory level
		parentDir := filepath.Dir(dir)
		if parentDir == dir {
			// Reached the root directory, stop searching
			break
		}

		dir = parentDir
	}

	return errors.New(".env file not found in the project hierarchy")
}

// setRootWorkingDirectoryEnvironment sets the root working directory environment.
func setRootWorkingDirectoryEnvironment() {
	absoluteRootWorkingDirectory := GetProjectRootWorkingDirectory()

	// when we `Set` a viper with string value, we should get it from viper with `viper.GetString`, elsewhere we get empty string
	viper.Set(constants.AppRootPath, absoluteRootWorkingDirectory)
}
