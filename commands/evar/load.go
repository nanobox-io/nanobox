package evar

import (
	"fmt"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/helpers"
	"github.com/nanobox-io/nanobox/models"
	app_evar "github.com/nanobox-io/nanobox/processors/app/evar"
	production_evar "github.com/nanobox-io/nanobox/processors/evar"
	"github.com/nanobox-io/nanobox/util"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

// LoadCmd loads variables from a file.
var LoadCmd = &cobra.Command{
	Use:   "load filename",
	Short: "Loads environment variable(s) from a file",
	Long:  ``,
	Run:   loadFn,
}

// loadFn parses a specified file and adds the contained variables to nanobox.
// Read in the file, strip out 'export ', parse, add resulting vars
func loadFn(ccmd *cobra.Command, args []string) {
	// parse the evars excluding the context
	env, _ := models.FindEnvByID(config.EnvID())
	args, location, name := helpers.Endpoint(env, args, 0)
	vars, err := loadVars(args)
	if err != nil {
		display.CommandErr(util.Err{
			Message: err.Error(),
			Code:    "USER",
			Stack:   []string{"failed to load evars from file"},
		})
		return
	}

	evars := parseEvars(vars)

	switch location {
	case "local":
		app, _ := models.FindAppBySlug(config.EnvID(), name)
		display.CommandErr(app_evar.Add(env, app, evars))
	case "production":
		steps.Run("login")(ccmd, args)

		production_evar.Add(env, name, evars)
	}
}

// loadVars loads variables from filenames passed in
func loadVars(args []string) ([]string, error) {
	vars := []string{}

	for _, filename := range args {
		contents, err := ioutil.ReadFile(filename)
		if err != nil {
			return nil, fmt.Errorf("Failed to read file - %s", err.Error())
		}

		// normalize file `key=val`
		newthings := strings.Replace(string(contents), "export ", "", -1)

		// strip out blank lines
		newthings = regexp.MustCompilePOSIX(`\n\n+`).ReplaceAllString(newthings, "\n")

		// strip trailing newline
		newthings = regexp.MustCompilePOSIX(`\n$`).ReplaceAllString(newthings, "")

		// get index of variable start (change regex to allow more than alphabet chars as keys)
		indexes := regexp.MustCompilePOSIX(`(^[a-zA-Z]*?)=(\"\n|.)`).FindAllStringIndex(newthings, -1)

		start := 0
		for i := range indexes {
			end := indexes[i][0]
			if end == 0 {
				continue
			}
			// end-1 leaves off the newline after the variable declaration
			vars = append(vars, newthings[start:end-1])
			start = end
		}
		// the newline after this variable declaration would have been previously stripped off
		vars = append(vars, newthings[start:])
	}

	return vars, nil
}
