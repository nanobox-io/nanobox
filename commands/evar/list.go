package evar

import (
	"strings"

	"github.com/jcelliott/lumber"
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
	Use:   "ls [local|dry-run|remote-alias]",
	Short: "list environment variable(s)",
	Long:  ``,
	Run:   listFn,
}

// listFn ...
func listFn(ccmd *cobra.Command, args []string) {

	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)

	// if the first argument is not a keyvalue pair, (at this point the
	// remote-alias would be stripped from helpers.Endpoint) it is likely
	// the app name. try setting vars on that app
	if len(args) > 0 && !strings.Contains(args[0], "=") {
		lumber.Info("Remote alias not found for '%s', attempting to set vars on app named '%s'\n", args[0], args[0])
		name = args[0]
		args = args[1:]
	}

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
