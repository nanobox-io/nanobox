package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"
)

// RunCmd ...
var RunCmd = &cobra.Command{
	Use:    "run",
	Short:  "Start your local development environment.",
	Long:   `
Starts your local development enviroment and opens an
interactive console inside the environment.

You can also pass a command into 'run'. Nanobox will
run the command without dropping you into a console
in your local environment.
	`,
	PreRun: steps.Run("start", "build-runtime", "dev start", "dev deploy"),
	Run:    runFn,
	PostRun: steps.Run("dev stop"),
}

// runFn ...
func runFn(ccmd *cobra.Command, args []string) {

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

	consoleConfig := console.ConsoleConfig{
		IsDev: true,
		DevIP: appModel.GlobalIPs["env"],
	}

	if len(args) > 0 {
		consoleConfig.Command = strings.Join(args, " ")
	}
	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(dev.Console(envModel, appModel, consoleConfig))
}
