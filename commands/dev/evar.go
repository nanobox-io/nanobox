package dev

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/util/config"
)

var (

	// EnvCmd ...
	EnvCmd = &cobra.Command{
		Use:   "evar",
		Short: "Manages environment variables in your local dev app.",
		Long:  ``,
	}

	// EnvAddCmd ...
	EnvAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds environment variable(s) to your dev app.",
		Long: `
Adds environment variable(s) to your dev app. Multiple key-value
pairs can be added simultaneously using a comma-delimited list.
		`,
		Run: envAddFn,
	}

	// EnvListCmd ...
	EnvListCmd = &cobra.Command{
		Use:   "ls",
		Short: "Lists all environment variables registered in your dev app.",
		Long:  ``,
		Run:   envListFn,
	}

	// EnvRemoveCmd ...
	EnvRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes environment variable(s) from your dev app.",
		Long: `
Removes environment variable(s) from your dev app. Multiple keys
can be removed simultaneously using a comma-delimited list.
		`,
		Run: envRemoveFn,
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
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
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
				app.Evars[strings.ToUpper(parts[0])] = parts[1]
			}
		}
	}
	fmt.Printf("app evars are update to: %+v\n", app.Evars)

	app.Save()
}

// envListFn ...
func envListFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	fmt.Println(app.Evars)
}

// envRemoveFn ...
func envRemoveFn(ccmd *cobra.Command, args []string) {
	app, _ := models.FindAppBySlug(config.EnvID(), "dev")
	for _, arg := range args {
		for _, key := range strings.Split(arg, ",") {
			delete(app.Evars, strings.ToUpper(key))
		}
	}

	app.Save()
}
