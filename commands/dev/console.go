package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

// ConsoleCmd ...
var ConsoleCmd = &cobra.Command{
	Use:    "console",
	Short:  "Opens an interactive console inside your dev platform.",
	Long:   ``,
	PreRun: steps.Run("start", "build-runtime", "dev start", "dev deploy"),
	Run:    consoleFn,
}

// consoleFn ...
func consoleFn(ccmd *cobra.Command, args []string) {

	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	// if given an argument they wanted to run a console into a container
	// if no arguement is provided they wanted to run a dev console
	// and be dropped into a dev environment
	if len(args) > 0 {
		component, _ := models.FindComponentBySlug(appModel.ID, args[0])

		display.CommandErr(env.Console(component, console.ConsoleConfig{}))
		return
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(dev.Console(envModel, appModel, false))
}
