package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/env"
)

var (

	// EnvCmd ...; currently a hidden command because its only used for one thing
	EnvCmd = &cobra.Command{
		Use:   "env",
		Short: "Shared environment provisioning",
		Long: `
A namespaced collection of hidded subcommands used primarily as share provisioning processes.
		`,
		Hidden: true,
	}
)

//
func init() {
	// hidden subcommands
	EnvCmd.AddCommand(env.ShareCmd)
}
