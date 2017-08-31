package evar

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	app_evar "github.com/nanobox-io/nanobox/processors/app/evar"
	production_evar "github.com/nanobox-io/nanobox/processors/evar"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// RemoveCmd removes an evar.
var RemoveCmd = &cobra.Command{
	Use:   "rm [local|dry-run|remote-alias] key",
	Short: "Remove environment variable(s)",
	Long: `Remove environment variable(s).

The alias must be used when removing variables from a production app.
If you would like to remove variables from a different app, please add
it as a remote as follows: 'nanobox remote add <APPNAME> <ALIAS>'.
You may then perform the 'rm' again, specifying that alias.`,
	Run: removeFn,
}

// removeFn ...
func removeFn(ccmd *cobra.Command, args []string) {
	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)

	if len(args) < 1 {
		fmt.Printf(`
--------------------------------------------
Please provide the key you'd like to remove!
--------------------------------------------

`)
		ccmd.HelpFunc()(ccmd, args)
		return
	}

	evars := parseKeys(args)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_evar.Remove(app, evars))
	case "production":
		steps.Run("login")(ccmd, args)

		env, _ := models.FindEnvByID(config.EnvID())

		production_evar.Remove(env, name, evars)
	}
}

func parseKeys(args []string) []string {
	keys := []string{}

	for _, arg := range args {
		// todo: remove support for comma separated key list (not necessary,
		// but would be nice to have uniform syntax not supporting it)
		for _, key := range strings.Split(arg, ",") {
			keys = append(keys, strings.ToUpper(key))
		}
	}

	return keys
}
