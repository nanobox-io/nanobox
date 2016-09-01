package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// StartCmd ...
	StartCmd = &cobra.Command{
		Use:   "start",
		Short: "Starts your dev platform.",
		Long: `
Starts your dev platform from its previous state. If starting for
the first time, you should also generate a build (nanobox build)
and deploy it into your dev platform (nanobox dev deploy).
		`,
		PreRun: validate.Requires("provider"),
		Run:    devStart,
	}
)

func init() {
	steps.Build("dev start", startCheck, devStart)
}

// devStart ...
func devStart(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")

	display.CommandErr(dev.Start(env, app))
}


func startCheck() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	return app.Status == "up"
}