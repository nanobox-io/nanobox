package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// DeployCmd ...
var DeployCmd = &cobra.Command{
	Use:    "deploy",
	Short:  "Deploys a build package into your dev platform and starts all data services.",
	Long:   ``,
	PreRun: steps.Run("start", "build", "compile", "dev start"),
	Run:    deployFn,
}

func init() {
	steps.Build("dev deploy", deployComplete, deployFn)
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(env.ID, "dev")
	display.CommandErr(dev.Deploy(env, app))
}

func deployComplete() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	env, _ := app.Env()
 	return app.DeployedBoxfile != "" && env.BuiltBoxfile == app.DeployedBoxfile
}
