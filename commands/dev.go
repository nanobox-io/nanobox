package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/dev"
)

var (

	// DevCmd ...
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Manages your 'development' environment.",
		Long:  ``,
		Hidden: true,
	}
)

//
func init() {

	// public subcommands
	DevCmd.AddCommand(dev.StartCmd)
	DevCmd.AddCommand(dev.StopCmd)
	DevCmd.AddCommand(dev.DeployCmd)
	DevCmd.AddCommand(dev.DestroyCmd)
}
