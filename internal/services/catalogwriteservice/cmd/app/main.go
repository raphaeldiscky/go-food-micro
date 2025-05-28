// Package main contains the main function for the catalogs write service.
package main

import (
	"os"

	"github.com/pterm/pterm"
	"github.com/pterm/pterm/putils"
	"github.com/spf13/cobra"

	"github.com/raphaeldiscky/go-food-micro/internal/services/catalogwriteservice/internal/shared/app"
)

func newRootCmd() *cobra.Command {
	return &cobra.Command{
		Use:              "catalogs-write-microservice",
		Short:            "catalogs-write-microservice based on vertical slice architecture",
		Long:             `This is a command runner or cli for api architecture in golang.`,
		TraverseChildren: true,
		Run: func(_ *cobra.Command, _ []string) {
			app.NewCatalogWriteApp().Run()
		},
	}
}

// https://github.com/swaggo/swag#how-to-use-it-with-gin

// @contact.name Raphael Discky
// @contact.url https://github.com/raphaeldiscky
// @title Catalogs Write-Service Api
// @version 1.0
// @description Catalogs Write-Service Api.
func main() {
	err := pterm.DefaultBigText.WithLetters(
		putils.LettersFromStringWithStyle("Catalogs", pterm.FgLightGreen.ToStyle()),
		putils.LettersFromStringWithStyle(" Write Service", pterm.FgLightMagenta.ToStyle())).
		Render()
	if err != nil {
		os.Exit(1)
	}

	err = newRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
