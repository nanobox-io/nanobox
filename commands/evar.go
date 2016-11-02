package commands

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/evar"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// EvarCmd ...
	EvarCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manages environment variables in your production app.",
		Long:  ``,
	}

	// EvarAddCmd ...
	EvarAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds environment variable(s) to your production app.",
		Long: `
Adds environment variable(s) to your production app. Multiple key-value
pairs can be added simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("login"),
		Run:    evarAddFn,
	}

	// EvarListCmd ...
	EvarListCmd = &cobra.Command{
		Use:    "ls",
		Short:  "Lists all environment variables registered in your production app.",
		Long:   ``,
		PreRun: steps.Run("login"),
		Run:    evarListFn,
	}

	// EvarRemoveCmd ...
	EvarRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes environment variable(s) from your production app.",
		Long: `
Removes environment variable(s) from your production app. Multiple keys
can be removed simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("login"),
		Run:    evarRemoveFn,
	}

	app string
)

//
func init() {
	EvarCmd.AddCommand(EvarAddCmd)
	EvarCmd.AddCommand(EvarRemoveCmd)
	EvarCmd.AddCommand(EvarListCmd)
	EvarCmd.PersistentFlags().StringVarP(&deployCmdFlags.app, "app", "a", "default", "message to accompany this command")
}

// evarAddFn ...
func evarAddFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	evars := map[string]string{}

	for _, arg := range args {
		// define a function that will allow us to
		// split on ',' or ' '
		f := func(c rune) bool {
			return c == ',' || c == ' '
		}

		for _, pair := range strings.FieldsFunc(arg, f) {
			// define a field split that llows us to split on
			// ':' or '='
			parts := strings.FieldsFunc(pair, func(c rune) bool {
				return c == ':' || c == '='
			})
			if len(parts) == 2 {

				evars[strings.ToUpper(parts[0])] = parts[1]
			}
		}
	}

	if app == "" {
		app = "default"
	}

	display.CommandErr(evar.Add(env, app, evars))
}

// evarListFn ...
func evarListFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	if app == "" {
		app = "default"
	}

	display.CommandErr(evar.List(env, app))
}

// evarRemoveFn ...
func evarRemoveFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	keys := []string{}

	for _, arg := range args {
		for _, key := range strings.Split(arg, ",") {
			keys = append(keys, key)
		}
	}

	if app == "" {
		app = "default"
	}

	display.CommandErr(evar.Remove(env, app, keys))
}