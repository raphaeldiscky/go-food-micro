// Package main contains the main function for the catalog read service.
package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogreadservice/internal/shared/app"
)

// newRootCmd creates and returns the root command for the catalog read service.
func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "catalogs-read-microservices",
		Short:            "catalogs-read-microservices based on vertical slice architecture",
		Long:             `This is a command runner or cli for api architecture in golang.`,
		TraverseChildren: true,
		Run: func(_ *cobra.Command, _ []string) {
			app.NewCatalogReadApp().Run()
		},
	}
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Raphael Discky
// @contact.url https://github.com/raphaeldiscky
// @title Catalogs Read-Service Api
// @version 1.0
// @description Catalogs Read-Service Api.
func main() {
	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Catalogs", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Read Service", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err != nil {
		os.Exit(1)
	}

	err = newRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
