package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/commands/steps"
)

// RunCmd ...
var RunCmd = &cobra.Command{
	Use:    "run",
	Short:  "Opens an dev container and starts all the code commands init.",
	Long:   ``,
	PreRun: steps.Run("start", "build", "dev start", "dev deploy"),
	Run:    runFn,
}

// runFn ...
func runFn(ccmd *cobra.Command, args []string) {

	// if given an argument they wanted to run a console into a container
	// if no arguement is provided they wanted to run a dev console
	// and be dropped into a dev environment
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(dev.Console(app, true))
}
