package evar

import (
	"fmt"
	"strconv"
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
	Use:   "add key=val[ key=val,key=val]",
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

// parseEvars parses evars already split into key="val" pairs.
func parseEvars(args []string) map[string]string {
	evars := map[string]string{}

	for _, pair := range args {
		// return after first split (in case there are `=` in the variable)
		parts := strings.SplitN(pair, "=", 2)
		if len(parts) == 2 {
			if parts[0] == "" {
				fmt.Printf(`
--------------------------------------------
Please provide a key to add evar!
--------------------------------------------`)
				continue
			}
			part := parts[1]
			var err error
			// if we've quoted the variable, and it's not a multiline, un-escape it
			if (len(parts[1]) > 1 && parts[1][0] == '"') && !strings.Contains(part, "\n") {
				// un-escape string values ("ensures proper escaped values too")
				// part, err = strconv.Unquote(strconv.Quote(parts[1]))
				part, err = strconv.Unquote(parts[1])
				if err != nil {
					fmt.Printf(`
--------------------------------------------
Please provide a properly escaped value!
--------------------------------------------`)
					continue
				}
			} else { // else, it's likely a multiline and we'll need to just remove quotes
				// strip var leading quote
				if parts[1][0] == '"' && len(parts[1]) > 1 {
					parts[1] = parts[1][1:]
				}

				// strip var ending quote
				if parts[1][len(parts[1])-1] == '"' && len(parts[1]) > 1 {
					parts[1] = parts[1][:len(parts[1])-1]
				}
				part = parts[1]
			}

			evars[strings.ToUpper(parts[0])] = part
		} else {
			fmt.Printf(`
--------------------------------------------
Please provide a valid evar! ("key=value")
--------------------------------------------`)
		}
	}

	return evars
}
