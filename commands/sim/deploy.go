package sim

import (
	"github.com/spf13/cobra"
	"github.com/nanobox-io/nanobox-boxfile"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/commands/registry"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (
	// DeployCmd ...
	DeployCmd = &cobra.Command{
		Use:   "deploy",
		Short: "Deploys a build package into your sim platform and starts all services.",
		Long: `
	Deploys a build package into your sim platform and
	starts all services. This is used to simulate a full
	deploy locally, before deploying into production.
	`,
		PreRun: func(ccmd *cobra.Command, args []string) {
			registry.Set("skip-compile", deployCmdFlags.skipCompile)
			steps.Run("start", "build-runtime", "compile-app", "sim start")(ccmd, args)
		},
		Run:    deployFn,
	}

	// deployCmdFlags ...
	deployCmdFlags = struct {
		skipCompile bool
	}{}
)


func init() {
	steps.Build("sim deploy", deployComplete, deployFn)
	DeployCmd.Flags().BoolVarP(&deployCmdFlags.skipCompile, "skip-compile", "", false, "skip compiling the app")
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(env.ID, "sim")
	// TODO: display an error if we cant find either of these

	display.CommandErr(sim.Deploy(env, app))
}

func deployComplete() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	env, _ := app.Env()
	return app.DeployedBoxfile != "" && env.BuiltBoxfile == app.DeployedBoxfile && buildComplete()
}

func buildComplete() bool {
	env, _ := models.FindEnvByID(config.EnvID())
	box := boxfile.NewFromPath(config.Boxfile())

	return env.UserBoxfile != "" && env.UserBoxfile == box.String()
}
