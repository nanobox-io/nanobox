package sim

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

var (

	// EnvCmd ...
	EnvCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manages environment variables in your local sim app.",
		Long:  ``,
	}

	// EnvAddCmd ...
	EnvAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds environment variable(s) to your sim app.",
		Long: `
Adds environment variable(s) to your sim app. Multiple key-value
pairs can be added simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("start", "build", "sim start", "sim deploy"),
		Run:    envAddFn,
	}

	// EnvListCmd ...
	EnvListCmd = &cobra.Command{
		Use:    "ls",
		Short:  "Lists all environment variables registered in your sim app.",
		Long:   ``,
		PreRun: steps.Run("start", "build", "sim start", "sim deploy"),
		Run:    envListFn,
	}

	// EnvRemoveCmd ...
	EnvRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes environment variable(s) from your sim app.",
		Long: `
Removes environment variable(s) from your sim app. Multiple keys
can be removed simultaneously using a comma-delimited list.
		`,
		PreRun: steps.Run("start", "build", "sim start", "sim deploy"),
		Run:    envRemoveFn,
	}
)

//
func init() {
	EnvCmd.AddCommand(EnvAddCmd)
	EnvCmd.AddCommand(EnvRemoveCmd)
	EnvCmd.AddCommand(EnvListCmd)
}

// envAddFn ...
func envAddFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	for _, arg := range args {
		for _, pair := range strings.Split(arg, ",") {
			parts := strings.FieldsFunc(pair, func(c rune) bool {
				return c == ':' || c == '='
			})
			if len(parts) == 2 {
				app.Evars[strings.ToUpper(parts[0])] = parts[1]
			}
		}
	}

	app.Save()
}

// envListFn ...
func envListFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	fmt.Println(app.Evars)
}

// envRemoveFn ...
func envRemoveFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "sim")
	for _, arg := range args {
		for _, key := range strings.Split(arg, ",") {
			delete(app.Evars, strings.ToUpper(key))
		}
	}

	app.Save()
}
