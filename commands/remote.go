package commands

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/steps"
	"github.com/nanobox-io/nanobox/models"
	"github.com/nanobox-io/nanobox/processors/remote"
	"github.com/nanobox-io/nanobox/util/config"
	"github.com/nanobox-io/nanobox/util/display"
)

var (

	// RemoteCmd ...
	RemoteCmd = &cobra.Command{
		Use:    "remote",
		Short:  "Manages remotes between local & production apps.",
		Long:   ``,
		PreRun: steps.Run("login"),
	}

	// RemoteAddCmd ...
	RemoteAddCmd = &cobra.Command{
		Use:   "add",
		Short: "Adds a new remote between a local & production app.",
		Long: `
Adds a new remote between a local and production app. A local
app can be remoted to multiple production apps. Each remote needs
an alias. If no alias is provided, 'default' is assumed.
		`,
		PreRun: steps.Run("login"),
		Run:    remoteAddFn,
	}

	// RemoteListCmd ...
	RemoteListCmd = &cobra.Command{
		Use:    "ls",
		Short:  "Lists all remotes for the current local app.",
		Long:   ``,
		PreRun: steps.Run("login"),
		Run:    remoteListFn,
	}

	// RemoteRemoveCmd ...
	RemoteRemoveCmd = &cobra.Command{
		Use:    "rm",
		Short:  "Removes a remote between a local & production app.",
		Long:   ``,
		PreRun: steps.Run("login"),
		Run:    remoteRemoveFn,
	}

)

//
func init() {
	RemoteCmd.AddCommand(RemoteAddCmd)
	RemoteCmd.AddCommand(RemoteListCmd)
	RemoteCmd.AddCommand(RemoteRemoveCmd)
}

// remoteAddFn ...
func remoteAddFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())

	if len(args) < 1 {
		fmt.Printf("\n! Please provide an app name to remote to\n\n")
		return
	}
	alias := ""
	if len(args) > 1 {
		alias = args[1]
	}
	display.CommandErr(remote.Add(env, args[0], alias))
}

// remoteListFn ...
func remoteListFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())

	display.CommandErr(remote.List(env))
}

// remoteRemoveFn ...
func remoteRemoveFn(ccmd *cobra.Command, args []string) {
	env, _ := models.FindEnvByID(config.EnvID())
	// set the meta arguments to be used in the processor and run the processor
	if len(args) == 0 {
		fmt.Println("invalid remote")
		return
	}

	display.CommandErr(remote.Remove(env, args[0]))
}
