package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts your sim platform.",
		Long: `
Starts the sim platform from its previous state. If starting for
the first time, you should also generate a build (nanobox build)
and deploy it into your sim platform (nanobox sim deploy).
		`,
		PreRun: steps.Run("start"),
		Run:    startFn,
	}
)

func init() {
	steps.Build("sim start", true, startCheck, startFn)
}

// startFn ...
func startFn(ccmd *cobra.Command, args []string) {
	// TODO: check the errors
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")

	display.CommandErr(sim.Start(env, app))
}

func startCheck() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	return app.Status == "up"
}
