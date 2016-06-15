package commands

import (
	"github.com/spf13/cobra"

	"github.com/nanobox-io/nanobox/commands/dev"
	"github.com/nanobox-io/nanobox/processor"
	"github.com/nanobox-io/nanobox/util/print"
	"github.com/nanobox-io/nanobox/validate"
)

var (

	// DevCmd ...
	DevCmd = &cobra.Command{
		Use:   "dev",
		Short: "Starts the Nanobox VM, provisions app, & opens an interactive terminal.",
		Long: `
Starts the Nanobox VM, provisions app, & opens an interactive
terminal. This is the primary command for managing local dev
apps and interacting with your Nanobox VM.
		`,
		PreRun: validate.Requires("provider"),
		Run:    devFn,
	}
)

//
func init() {

	// hidden subcommands
	DevCmd.AddCommand(dev.NetfsCmd)

	// public subcommands
	DevCmd.AddCommand(dev.DeployCmd)
	DevCmd.AddCommand(dev.DestroyCmd)
	DevCmd.AddCommand(dev.DNSCmd)
	DevCmd.AddCommand(dev.InfoCmd)
	DevCmd.AddCommand(dev.ExecCmd)
	DevCmd.AddCommand(dev.ConsoleCmd)
	DevCmd.AddCommand(dev.EnvCmd)
	DevCmd.AddCommand(dev.ResetCmd)
}

// devFn ...
func devFn(ccmd *cobra.Command, args []string) {
	print.OutputCommandErr(processor.Run("dev", processor.DefaultConfig))
}
