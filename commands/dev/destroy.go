package dev

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
	"github.com/nanobox-io/nanobox/processor/dev"
	"github.com/nanobox-io/nanobox/util/config"
)

// DestroyCmd ...
var (
	DestroyCmd = &cobra.Command{
		Use:    "destroy",
		Short:  "Destroys the docker machines associated with your dev app.",
		Long:   ``,
		PreRun: validate.Requires("provider", "provider_up"),
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
	devDestroy := dev.Destroy{App: getApp()}

	print.OutputCommandErr(devDestroy.Run())
}

// look up the real app id based on what they told me.
func getApp() models.App {
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
