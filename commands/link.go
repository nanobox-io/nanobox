package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/link"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
	"github.com/nanobox-io/nanobox/commands/steps"
)

var (

	// LinkCmd ...
	LinkCmd = &cobra.Command{
		Use:   "link",
		Short: "Manages links between local & production apps.",
		Long:  ``,
		PreRun: steps.Run("login"),
		Run:   linkAddFn,
	}

	// LinkAddCmd ...
	LinkAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new link between a local & production app.",
		Long: `
Adds a new link between a local and production app. A local
app can be linked to multiple production apps. Each link needs
an alias. If no alias is provided, 'default' is assumed.
		`,
		PreRun: steps.Run("login"),
		Run: linkAddFn,
	}

	// LinkListCmd ...
	LinkListCmd = &cobra.Command{
		Use:   "ls",
		Short: "Lists all links for the current local app.",
		Long:  ``,
		PreRun: steps.Run("login"),
		Run:   linkListFn,
	}

	// LinkRemoveCmd ...
	LinkRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a link between a local & production app.",
		Long:  ``,
		PreRun: steps.Run("login"),
		Run:   linkRemoveFn,
	}

	// linkCmdFlags ...
	linkCmdFlags = struct {
		alias string
	}{}
)

//
func init() {
	LinkCmd.PersistentFlags().StringVarP(&linkCmdFlags.alias, "alias", "a", "", "alias")

	LinkCmd.AddCommand(LinkAddCmd)
	LinkCmd.AddCommand(LinkListCmd)
	LinkCmd.AddCommand(LinkRemoveCmd)
}

// linkAddFn ...
func linkAddFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	// TODO: validate env?
	if len(args) != 1 {
		fmt.Printf("TODO: message: i need an app name to link to")
	}

	display.CommandErr(link.Add(env, args[0], linkCmdFlags.alias))
}

// linkListFn ...
func linkListFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())

	display.CommandErr(link.List(env))
}

// linkRemoveFn ...
func linkRemoveFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	// set the meta arguments to be used in the processor and run the processor
	if len(args) != 0 {
		linkCmdFlags.alias = args[0]
	}

	display.CommandErr(link.Remove(env, linkCmdFlags.alias))
}
