package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

func init() {
	steps.Build("dev start", true, startCheck, devStart)
}

// devStart ...
func devStart(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	appModel, _ := models.FindAppBySlug(config.EnvID(), "dev")
	
	display.CommandErr(env.Setup(envModel))
	display.CommandErr(app.Start(envModel, appModel, "dev"))
}

func startCheck() bool {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	return app.Status == "up"
}
