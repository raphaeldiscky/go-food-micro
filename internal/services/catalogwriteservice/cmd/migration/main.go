package main

import (
	"context"
	"os"

	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/config/environment"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/external/fxlog"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/zap"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/contracts"
	"github.com/raphaeldiscky/go-food-micro/internal/pkg/migration/goose"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	defaultLogger "github.com/raphaeldiscky/go-food-micro/internal/pkg/logger/defaultlogger"
	gormPostgres "github.com/raphaeldiscky/go-food-micro/internal/pkg/postgresgorm"

	appconfig "github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/config"
)

// newRootCmd creates and returns the root command for migrations.
func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
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

	// Add commands to the root command
	cmd.AddCommand(newDownCmd())
	cmd.AddCommand(newUpCmd())

	return cmd
}

// newDownCmd creates and returns the down migration command.
func newDownCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "down",
		Short: "Run a down migration",
		Run: func(cmd *cobra.Command, _ []string) {
			executeMigration(cmd, migration.Down)
		},
	}

	// Add flags to specify the version
	cmd.Flags().Uint("version", 0, "Migration version")

	return cmd
}

// newUpCmd creates and returns the up migration command.
func newUpCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "up",
		Short: "Run an up migration",
		Run: func(cmd *cobra.Command, _ []string) {
			executeMigration(cmd, migration.Up)
		},
	}

	// Add flags to specify the version
	cmd.Flags().Uint("version", 0, "Migration version")

	return cmd
}

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
		appconfig.NewModule(),
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
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		defaultLogger.GetLogger().Error(err)
		os.Exit(1)
	}
}
