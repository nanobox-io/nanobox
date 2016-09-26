package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/util/config"
)

var (

	// DestroyCmd ...
	DestroyCmd = &cobra.Command{
		Use:    "destroy",
		Short:  "Destroys the current project and removes it from Nanobox",
		Long:   `
Destroys the current project and removes it from Nanobox, destroying
the filesystem mount, associated dns aliases, and local app data.
		`,
		PreRun: steps.Run("start"),
		Run:    destroyFunc,
	}
)

// destroyFunc ...
func destroyFunc(ccmd *cobra.Command, args []string) {
	envModel, _ := models.FindEnvByID(config.EnvID())
	display.CommandErr(env.Destroy(envModel))
}
