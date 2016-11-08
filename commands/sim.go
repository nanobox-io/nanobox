package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/sim"
)

var (

	// SimCmd ...
	SimCmd = &cobra.Command{
		Use:   "sim",
		Short: "Manages your 'simulated' environment.",
		Long:  ``,
		Hidden: true,
	}
)

//
func init() {

	// public subcommands
	SimCmd.AddCommand(sim.StartCmd)
	SimCmd.AddCommand(sim.StopCmd)
	SimCmd.AddCommand(sim.DeployCmd)
	SimCmd.AddCommand(sim.DestroyCmd)
}
