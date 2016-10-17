package sim

import (
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/app/evar"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// EvarCmd ...
	EvarCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manages environment variables in your local sim app.",
		Long:  ``,
	}

	// EvarAddCmd ...
	EvarAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds environment variable(s) to your sim app.",
		Long: `
Adds environment variable(s) to your sim app. Multiple key-value
pairs can be added simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("start", "build", "compile", "sim start", "sim deploy"),
		Run:    evarAddFn,
	}

	// EvarListCmd ...
	EvarListCmd = &cobra.Command{
		Use:    "ls",
		Short:  "Lists all environment variables registered in your sim app.",
		Long:   ``,
		PreRun: steps.Run("start", "build", "compile", "sim start", "sim deploy"),
		Run:    evarListFn,
	}

	// EvarRemoveCmd ...
	EvarRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes environment variable(s) from your sim app.",
		Long: `
Removes environment variable(s) from your sim app. Multiple keys
can be removed simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("start", "build", "compile", "sim start", "sim deploy"),
		Run:    evarRemoveFn,
	}
)

//
func init() {
	EvarCmd.AddCommand(EvarAddCmd)
	EvarCmd.AddCommand(EvarRemoveCmd)
	EvarCmd.AddCommand(EvarListCmd)
}

// evarAddFn ...
func evarAddFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
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

	display.CommandErr(evar.Add(app, evars))
}

// evarListFn ...
func evarListFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	display.CommandErr(evar.List(app))
}

// evarRemoveFn ...
func evarRemoveFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	keys := []string{}

	for _, arg := range args {
		for _, key := range strings.Split(arg, ",") {
			keys = append(keys, strings.ToUpper(key))
		}
	}

	display.CommandErr(evar.Remove(app, keys))
}
