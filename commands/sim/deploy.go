package sim

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/sim"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

// DeployCmd ...
var DeployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploys a build package into your sim platform and starts all services.",
	Long: `
Deploys a build package into your sim platform and
starts all services. This is used to simulate a full
deploy locally, before deploying into production.
		`,
	PreRun: validate.Requires("provider", "provider_up", "built"),
	Run:    deployFn,
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(env.ID, "sim")
	// TODO: display an error if we cant find either of these
	
	display.CommandErr(sim.Deploy(env, app))
}
