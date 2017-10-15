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
		Use:   "remote",
		Short: "Manage application remotes.",
		Long: `
Manages connections between your local codebase and
remote, live applications created with Nanobox.
		`,
		PreRun: steps.Run("login"),
	}

	// RemoteAddCmd ...
	RemoteAddCmd = &cobra.Command{
		Use:   "add <team-name>/<app-name> [remote-alias]",
		Short: "Add a new remote.",
		Long: `
Adds a new remote. A local app can have multiple remotes. Each
remote needs an alias. If no alias is provided, 'default' is assumed.
		`,
		PreRun: steps.Run("login"),
		Run:    remoteAddFn,
	}

	// RemoteListCmd ...
	RemoteListCmd = &cobra.Command{
		Use:    "ls",
		Short:  "List all remotes for the current local app.",
		Long:   ``,
		PreRun: steps.Run("login"),
		Run:    remoteListFn,
	}

	// RemoteRemoveCmd ...
	RemoteRemoveCmd = &cobra.Command{
		Use:    "rm [remote-alias]",
		Short:  "Remove a remote.",
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
		fmt.Printf("\n! Please provide the app name for your remote\n\n")
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
