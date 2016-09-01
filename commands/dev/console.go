package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/commands/steps"
)

// ConsoleCmd ...
var ConsoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Opens an interactive console inside your dev platform.",
	Long:   ``,
	PreRun: steps.Run("start", "build", "dev start", "dev deploy"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {

	app, _ := models.FindAppBySlug(config.EnvID(), "dev")

	// if given an argument they wanted to run a console into a container
	// if no arguement is provided they wanted to run a dev console
	// and be dropped into a dev environment
	if len(args) > 0 {
		component, _ := models.FindComponentBySlug(app.ID, args[0])

		display.CommandErr(env.Console(component, env.ConsoleConfig{}))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(dev.Console(app, false))
}
