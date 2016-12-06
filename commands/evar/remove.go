package evar

import (
	// "fmt"
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

// RemoveCmd ...
var RemoveCmd = &cobra.Command{
	Use:   "rm",
	Short: "Remove environment variable(s)",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: removeFn,
}

// removeFn ...
func removeFn(ccmd *cobra.Command, args []string) {
	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)
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
		for _, key := range strings.Split(arg, ",") {
			if key != "" {
				keys = append(keys, key)
			}
		}
	}

	return keys
}
