package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/console"
	"github.com/nanobox-io/nanobox/util/display"

	// imported because we need its steps added
	_ "github.com/nanobox-io/nanobox/commands/dev"
)

// RunCmd ...
var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Start your local development environment.",
	Long: `
Starts your local development enviroment and opens an
interactive console inside the environment.

You can also pass a command into 'run'. Nanobox will
run the command without dropping you into a console
in your local environment.
	`,
	PreRun:  steps.Run("configure", "start", "build-runtime", "dev start", "dev deploy"),
	Run:     runFn,
	PostRun: steps.Run("dev stop"),
}

// runFn ...
func runFn(ccmd *cobra.Command, args []string) {

	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")

	consoleConfig := console.ConsoleConfig{
		IsDev: true,
		DevIP: appModel.GlobalIPs["env"],
	}

	if len(args) > 0 {
		consoleConfig.Command = strings.Join(args, " ")
	}

	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(processors.Run(envModel, appModel, consoleConfig))
}

func init() {
	steps.Build("dev deploy", true, devDeployComplete, devDeploy)
}

// devDeploy ...
func devDeploy(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(envModel.ID, "dev")
	display.CommandErr(app.Deploy(envModel, appModel))
}

func devDeployComplete() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	env, _ := app.Env()
	return app.DeployedBoxfile != "" && env.BuiltBoxfile == app.DeployedBoxfile && buildComplete()
}
