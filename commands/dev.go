package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/dev"
)

var (

	// DevCmd ...
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "The development environment ",
		Long: `
The development environment all starts below this subcommand.
		`,
	}
)

//
func init() {

	// hidden subcommands
	DevCmd.AddCommand(dev.NetfsCmd)

	// public subcommands
	DevCmd.AddCommand(dev.StartCmd)
	DevCmd.AddCommand(dev.DeployCmd)
	DevCmd.AddCommand(dev.DestroyCmd)
	DevCmd.AddCommand(dev.DNSCmd)
	DevCmd.AddCommand(dev.InfoCmd)
	DevCmd.AddCommand(dev.ExecCmd)
	DevCmd.AddCommand(dev.ConsoleCmd)
	DevCmd.AddCommand(dev.EnvCmd)
	DevCmd.AddCommand(dev.ResetCmd)
}
