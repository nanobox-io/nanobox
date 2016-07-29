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
	}
)

//
func init() {

	// public subcommands
	DevCmd.AddCommand(dev.StartCmd)
	DevCmd.AddCommand(dev.StopCmd)
	DevCmd.AddCommand(dev.DeployCmd)
	DevCmd.AddCommand(dev.DestroyCmd)
	DevCmd.AddCommand(dev.DNSCmd)
	DevCmd.AddCommand(dev.InfoCmd)
	DevCmd.AddCommand(dev.LogCmd)
	DevCmd.AddCommand(dev.ConsoleCmd)
	DevCmd.AddCommand(dev.EnvCmd)
	DevCmd.AddCommand(dev.ResetCmd)
	DevCmd.AddCommand(dev.UpCmd)
}
