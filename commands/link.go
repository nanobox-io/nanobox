package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processor/link"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// LinkCmd ...
	LinkCmd = &cobra.Command{
		Use:   "link",
		Short: "Manages links between local & production apps.",
		Long:  ``,
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
		Run: linkAddFn,
	}

	// LinkListCmd ...
	LinkListCmd = &cobra.Command{
		Use:   "ls",
		Short: "Lists all links for the current local app.",
		Long:  ``,
		Run:   linkListFn,
	}

	// LinkRemoveCmd ...
	LinkRemoveCmd = &cobra.Command{
		Use:   "rm",
		Short: "Removes a link between a local & production app.",
		Long:  ``,
		Run:   linkRemoveFn,
	}

	// linkCmdFlags ...
	linkCmdFlags = struct {
		app   string
		alias string
	}{}
)

//
func init() {
	LinkAddCmd.Flags().StringVarP(&linkCmdFlags.app, "app", "n", "", "app name")
	LinkCmd.PersistentFlags().StringVarP(&linkCmdFlags.alias, "alias", "a", "", "alias")

	LinkCmd.AddCommand(LinkAddCmd)
	LinkCmd.AddCommand(LinkListCmd)
	LinkCmd.AddCommand(LinkRemoveCmd)
}

// linkAddFn ...
func linkAddFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	add := link.Add{
		Env:   env,
		App:   linkCmdFlags.app,
		Alias: linkCmdFlags.alias,
	}
	display.CommandErr(add.Run())
}

// linkListFn ...
func linkListFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	list := link.List{
		Env: env,
	}
	display.CommandErr(list.Run())
}

// linkRemoveFn ...
func linkRemoveFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	remove := link.Remove{
		Env:   env,
		Alias: linkCmdFlags.alias,
	}
	// set the meta arguments to be used in the processor and run the processor
	display.CommandErr(remove.Run())
}
