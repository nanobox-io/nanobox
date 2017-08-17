package evar

import (
	"fmt"
	"io/ioutil"
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

	for _, arg := range args {
		// define a function that will allow us to
		// split on ',' or ' '
		f := func(c rune) bool {
			return c == ','
		}

		for _, pair := range strings.FieldsFunc(arg, f) {
			// define a field split that llows us to split on
			// ':' or '='

			// todo: return after first split (in case there are `=` in the variable)
			// parts := strings.FieldsFunc(pair, func(c rune) bool {
			// 	return c == '='
			// })

			parts := strings.SplitN(pair, "=", 2)
			if len(parts) == 2 {
				// check to see if the value is a file
				content, err := ioutil.ReadFile(parts[1])
				if err == nil {
					parts[1] = string(content)
				}
				evars[strings.ToUpper(parts[0])] = parts[1]
			} else {
				fmt.Printf("invalid evar (%s)\n", pair)
			}
		}
	}

	return evars
}
