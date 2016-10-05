package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// LogCmd ...
	LogCmd = &cobra.Command{
		Use:    "log",
		Short:  "Displays logs from the running dev app and its components.",
		Long:   ``,
		PreRun: steps.Run("start", "build", "compile", "dev start", "dev deploy"),
		Run:    logFn,
	}
)

// logFn will run the DNS processor for adding DNS entires to the "hosts" file
func logFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	display.CommandErr(dev.Log(app))
}
