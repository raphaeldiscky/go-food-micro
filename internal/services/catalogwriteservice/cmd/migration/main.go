package main

import (
	"context"
	"os"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/goose"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"
	appconfig "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"

	"github.com/spf13/cobra"
	"go.uber.org/fx"
)

func init() {
	// Add flags to specify the version
	cmdUp.Flags().Uint("version", 0, "Migration version")
	cmdDown.Flags().Uint("version", 0, "Migration version")

	// Add commands to the root command
	rootCmd.AddCommand(cmdUp)
	rootCmd.AddCommand(cmdDown)
}

var (
	rootCmd = &cobra.Command{
		Use:   "migration",
		Short: "A tool for running migrations",
		Run: func(cmd *cobra.Command, args []string) {
			// Execute the "up" subcommand when no subcommand is specified
			if len(args) == 0 {
				cmd.SetArgs([]string{"up"})
				if err := cmd.Execute(); err != nil {
					defaultLogger.GetLogger().Error(err)
					os.Exit(1)
				}
			}
		},
	}

	cmdDown = &cobra.Command{
		Use:   "down",
		Short: "Run a down migration",
		Run: func(cmd *cobra.Command, _ []string) {
			executeMigration(cmd, migration.Down)
		},
	}

	cmdUp = &cobra.Command{
		Use:   "up",
		Short: "Run an up migration",
		Run: func(cmd *cobra.Command, _ []string) {
			executeMigration(cmd, migration.Up)
		},
	}
)

func executeMigration(cmd *cobra.Command, commandType migration.CommandType) {
	version, err := cmd.Flags().GetUint("version")
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}

	app := fx.New(
		config.ModuleFunc(environment.Development),
		zap.Module,
		fxlog.FxLogger,
		gormPostgres.Module,
		appconfig.Module,
		// use go-migrate library for migration
		// gomigrate.Module,
		// use go-migrate library for migration
		goose.Module,
		fx.Invoke(
			func(migrationRunner contracts.PostgresMigrationRunner, logger logger.Logger) {
				logger.Info("Migration process started...")

				switch commandType {
				case migration.Up:
					err = migrationRunner.Up(context.Background(), version)
				case migration.Down:
					err = migrationRunner.Down(context.Background(), version)
				}

				if err != nil {
					logger.Fatalf("migration failed, err: %s", err)
				}

				logger.Info("Migration completed...")
			},
		),
	)

	err = app.Start(context.Background())
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}

	err = app.Stop(context.Background())
	if err != nil {
		defaultLogger.GetLogger().Fatal(err)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		defaultLogger.GetLogger().Error(err)
		os.Exit(1)
	}
}
