package evar

import (
	// "fmt"
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	app_evar "github.com/nanobox-io/nanobox/processors/app/evar"
	production_evar "github.com/nanobox-io/nanobox/processors/evar"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// ListCmd ...
var ListCmd = &cobra.Command{
	Use:   "ls",
	Short: "list environment variable(s)",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: listFn,
}

// listFn ...
func listFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_evar.List(app))
	case "production":
		steps.Run("login")(ccmd, args)

		env, _ := models.FindEnvByID(config.EnvID())
		production_evar.List(env, name)
	}
}
