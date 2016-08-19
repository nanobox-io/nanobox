package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

// DeployCmd ...
var DeployCmd = &cobra.Command{
	Use:    "deploy",
	Short:  "Deploys a build package into your dev platform and starts all data services.",
	Long:   ``,
	PreRun: validate.Requires("provider", "provider_up", "built", "dev_isup"),
	Run:    deployFn,
}

// deployFn ...
func deployFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(env.ID, "dev")
	devDeploy := dev.Deploy{
		Env: env,
		App: app,
	}
	display.CommandErr(devDeploy.Run())
}
