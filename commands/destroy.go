package commands

import (
	"fmt"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app"
	"github.com/nanobox-io/nanobox/processors/env"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// DestroyCmd ...
	DestroyCmd = &cobra.Command{
		Use:   "destroy",
		Short: "Destroy the current project and remove it from Nanobox.",
		Long: `
Destroys the current project and removes it from Nanobox – destroying
the filesystem mount, associated dns aliases, and local app data.
		`,
		PreRun: steps.Run("start"),
		Run:    destroyFunc,
	}
)

// destroyFunc ...
func destroyFunc(ccmd *cobra.Command, args []string) {
	envModel, err := models.FindEnvByID(config.EnvID())
	if err != nil {
		fmt.Println("This project doesn't exist on nanobox.")
		return
	}

	if len(args) == 0 {
		display.CommandErr(env.Destroy(envModel))
		return
	}

	_, _, name := helpers.Endpoint(envModel, args, 2)
	appModel, err := models.FindAppBySlug(envModel.ID, name)
	if err != nil {
		fmt.Println("Could not find the application")
	}

	display.CommandErr(app.Destroy(appModel))

}
