package main

import (
	"os"

	"github.com/raphaeldiscky/go-food-micro/internal/services/orderservice/internal/shared/app"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"
)

// newRootCmd creates and returns the root command for the orders service.
func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "orders-microservice",
		Short:            "orders-microservice based on vertical slice architecture",
		Long:             `This is a command runner or cli for api architecture in golang.`,
		TraverseChildren: true,
		Run: func(_ *cobra.Command, _ []string) {
			app.NewApp().Run()
		},
	}
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Raphael Discky
// @contact.url https://github.com/raphaeldiscky
// @title Orders Service Api
// @version 1.0
// @description Orders Service Api
func main() {
	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Orders", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Service", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err != nil {
		os.Exit(1)
	}

	rootCmd := newRootCmd()
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
