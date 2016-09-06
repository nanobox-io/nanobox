package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/dev"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// DestroyCmd ...
var (
	DestroyCmd = &cobra.Command{
		Use:    "destroy",
		Short:  "Destroys the docker machines associated with your dev app.",
		Long:   ``,
		PreRun: steps.Run("start", "dev start"),
		Run:    destroyFn,
	}

	destroyCmdFlags = struct {
		app string
	}{}
)

//
func init() {
	DestroyCmd.Flags().StringVarP(&destroyCmdFlags.app, "app", "a", "", "app to destroy")
}

// destroyFn ...
func destroyFn(ccmd *cobra.Command, args []string) {
	display.CommandErr(dev.Destroy(getApp()))
}

// look up the real app id based on what they told me.
func getApp() *models.App {
	if destroyCmdFlags.app != "" {
		envs, _ := models.AllEnvs()
		for _, env := range envs {
			app, _ := models.FindAppBySlug(env.ID, "dev")
			if env.ID == destroyCmdFlags.app || app.ID == destroyCmdFlags.app {
				return app
			}
		}
	}

	// if none could be found based on the arguements
	// use the one based on my folder
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	return app
}
