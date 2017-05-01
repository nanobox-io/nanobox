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

// AddCmd ...
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "Adds environment variable(s)",
	Long:  ``,
	// PreRun: steps.Run("login"),
	Run: addFn,
}

var literal bool

func init() {
	AddCmd.Flags().BoolVarP(&literal, "literal", "", false, "take the added evar at face value (no interpolation) only one key per command")
}
// addFn ...
func addFn(ccmd *cobra.Command, args []string) {

	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)
	evars := parseEvars(args)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_evar.Add(env, app, evars))
	case "production":
		steps.Run("login")(ccmd, args)

		production_evar.Add(env, name, evars)
	}
}

func parseEvars(args []string) map[string]string {
	evars := map[string]string{}

	if literal {
		for _, arg := range args {
			parts := strings.Split(arg, "=")
			if len(parts) < 2 {
				fmt.Printf("invalid evar (%s)\n", arg)
				return evars
			}
			key := parts[0]
			val := strings.Replace(arg, fmt.Sprintf("%s=", key), "", 1)
			evars[key] = val
			return evars
		}
		
	}


	for _, arg := range args {
		// define a function that will allow us to
		// split on ',' or ' '
		f := func(c rune) bool {
			return c == ','
		}

		for _, pair := range strings.FieldsFunc(arg, f) {
			// define a field split that allows us to split on
			// ':' or '='
			parts := strings.FieldsFunc(pair, func(c rune) bool {
				return c == '='
			})
			if len(parts) == 2 {

				evars[strings.ToUpper(parts[0])] = parts[1]
			} else {
				fmt.Printf("invalid evar (%s)\n", pair)
			}
		}
	}

	return evars
}
