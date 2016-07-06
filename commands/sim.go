package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/sim"
)

var (

	// SimCmd ...
	SimCmd = &cobra.Command{
		Use:   "sim",
		Short: "The simulated environment ",
		Long: `
The simulated environment all starts below this subcommand.
		`,
	}
)

//
func init() {

	// public subcommands
	SimCmd.AddCommand(sim.StartCmd)
	SimCmd.AddCommand(sim.StopCmd)
	SimCmd.AddCommand(sim.DeployCmd)
	SimCmd.AddCommand(sim.DestroyCmd)
	SimCmd.AddCommand(sim.InfoCmd)
	SimCmd.AddCommand(sim.LogCmd)
	SimCmd.AddCommand(sim.ConsoleCmd)
	SimCmd.AddCommand(sim.EnvCmd)
	SimCmd.AddCommand(sim.DNSCmd)
	SimCmd.AddCommand(sim.UpCmd)
}
